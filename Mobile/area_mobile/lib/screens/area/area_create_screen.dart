import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';

import '../../models/area_backend_models.dart';
import '../../models/area_definitions.dart';
import '../../providers/area_provider.dart';
import '../../providers/auth_provider.dart';
import '../../services/service_catalog_service.dart';
import '../../services/service_connector.dart';
import '../../theme/colors.dart';

enum AreaWizardStep { action, reaction, details }

class CreateAreaScreen extends StatefulWidget {
  const CreateAreaScreen({super.key});

  @override
  State<CreateAreaScreen> createState() => _CreateAreaScreenState();
}

class _CreateAreaScreenState extends State<CreateAreaScreen> {
  late final ServiceCatalogService _catalog;
  late final ServiceConnector _connector;

  bool _isLoading = true;
  String? _error;

  List<_MobileService> _services = [];

  _MobileService? _actionService;
  _MobileService? _reactionService;

  _MobileAction? _selectedAction;
  _MobileReaction? _selectedReaction;

  Map<String, String> _actionFieldValues = {};
  Map<String, String> _reactionFieldValues = {};

  final Map<String, TextEditingController> _reactionControllers = {};
  final Map<String, FocusNode> _reactionFocusNodes = {};
  String? _focusedReactionField;

  String _areaName = '';
  AreaWizardStep _wizardStep = AreaWizardStep.action;
  String? _createError;
  bool _isCreating = false;

  @override
  void initState() {
    super.initState();
    _catalog = ServiceCatalogService();
    _connector = ServiceConnector();
    _loadFromBackend();
  }

  @override
  void dispose() {
    for (final c in _reactionControllers.values) {
      c.dispose();
    }
    for (final f in _reactionFocusNodes.values) {
      f.dispose();
    }
    super.dispose();
  }

  Future<void> _loadFromBackend() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      final authProvider = context.read<AuthProvider>();
      final rawUserId = authProvider.user?['id'];
      final userId = rawUserId is int
          ? rawUserId
          : int.tryParse(rawUserId?.toString() ?? '');

      final Map<String, bool> connectedByProvider = {};
      if (userId != null) {
        final statuses = await _connector.fetchServices(userId);
        for (final status in statuses) {
          connectedByProvider[status.name] = status.isConnected;
        }
      }

      final serviceNames = await _catalog.fetchServiceNames();

      final configs = await Future.wait(
        serviceNames.map((name) async {
          try {
            return await _catalog.fetchServiceConfig(name);
          } catch (_) {
            return null;
          }
        }),
      );

      final built = configs.whereType<ServiceConfigDto>().map((cfg) {
        final isConnected = cfg.provider.trim().isEmpty
            ? true
            : (connectedByProvider[cfg.provider] ?? false);

        return _MobileService.fromConfig(cfg, connected: isConnected);
      }).toList();

      setState(() {
        _services = built;
      });
    } catch (e) {
      setState(() {
        _error = e.toString().replaceAll('Exception: ', '');
        _services = [];
      });
    } finally {
      setState(() {
        _isLoading = false;
      });
    }
  }

  Map<String, String> _initValues(List<AreaFieldDefinition> fields) {
    final out = <String, String>{};
    for (final f in fields) {
      out[f.name] = f.defaultValue ?? '';
    }
    return out;
  }

  bool _requiredFilled(
    List<AreaFieldDefinition> fields,
    Map<String, String> values,
  ) {
    return fields.every((f) {
      if (!f.required) return true;
      final v = values[f.name];
      return v != null && v.trim().isNotEmpty;
    });
  }

  bool get _canProceedAction {
    return _actionService != null &&
        _actionService!.connected &&
        _selectedAction != null &&
        _requiredFilled(_selectedAction!.fields, _actionFieldValues);
  }

  bool get _canProceedReaction {
    return _reactionService != null &&
        _reactionService!.connected &&
        _selectedReaction != null &&
        _requiredFilled(_selectedReaction!.fields, _reactionFieldValues);
  }

  bool get _canCreate {
    return _areaName.trim().isNotEmpty &&
        _canProceedAction &&
        _canProceedReaction;
  }

  void _next() {
    if (_wizardStep == AreaWizardStep.action && _canProceedAction) {
      setState(() {
        _wizardStep = AreaWizardStep.reaction;
        _createError = null;
      });
      return;
    }
    if (_wizardStep == AreaWizardStep.reaction && _canProceedReaction) {
      setState(() {
        _wizardStep = AreaWizardStep.details;
        _createError = null;
      });
    }
  }

  void _prev() {
    if (_wizardStep == AreaWizardStep.details) {
      setState(() {
        _wizardStep = AreaWizardStep.reaction;
        _createError = null;
      });
      return;
    }
    if (_wizardStep == AreaWizardStep.reaction) {
      setState(() {
        _wizardStep = AreaWizardStep.action;
        _createError = null;
      });
    }
  }

  Future<void> _pickDateTime(
    String fieldName,
    void Function(String, String) onChanged,
  ) async {
    final now = DateTime.now();
    final pickedDate = await showDatePicker(
      context: context,
      initialDate: now,
      firstDate: now.subtract(const Duration(days: 1)),
      lastDate: DateTime(2100),
    );
    if (pickedDate == null) return;

    final pickedTime = await showTimePicker(
      context: context,
      initialTime: TimeOfDay.now(),
    );
    if (pickedTime == null) return;

    final dt = DateTime(
      pickedDate.year,
      pickedDate.month,
      pickedDate.day,
      pickedTime.hour,
      pickedTime.minute,
    );

    onChanged(fieldName, dt.toUtc().toIso8601String());
  }

  void _initReactionEditors(List<AreaFieldDefinition> fields) {
    for (final c in _reactionControllers.values) {
      c.dispose();
    }
    for (final f in _reactionFocusNodes.values) {
      f.dispose();
    }
    _reactionControllers.clear();
    _reactionFocusNodes.clear();
    _focusedReactionField = null;

    for (final field in fields) {
      final initial =
          _reactionFieldValues[field.name] ?? (field.defaultValue ?? '');
      final controller = TextEditingController(text: initial);
      final focusNode = FocusNode();

      focusNode.addListener(() {
        if (focusNode.hasFocus) {
          setState(() {
            _focusedReactionField = field.name;
          });
        }
      });

      _reactionControllers[field.name] = controller;
      _reactionFocusNodes[field.name] = focusNode;
    }
  }

  Future<void> _insertOrCopyOutputPlaceholder(String outputName) async {
    final placeholder = '{{${outputName.trim()}}}';

    // Insert into last focused reaction field (even if focus was lost when tapping)
    final target = _focusedReactionField;
    if (target == null || !_reactionControllers.containsKey(target)) {
      await Clipboard.setData(ClipboardData(text: placeholder));
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Copié: $placeholder'),
          backgroundColor: context.appColors.deepBlue,
        ),
      );
      return;
    }

    final controller = _reactionControllers[target]!;
    final focusNode = _reactionFocusNodes[target];

    final text = controller.text;
    final sel = controller.selection;

    final start = (sel.start >= 0 && sel.start <= text.length)
        ? sel.start
        : text.length;
    final end =
        (sel.end >= 0 && sel.end <= text.length) ? sel.end : text.length;

    final newText = text.replaceRange(start, end, placeholder);
    controller.text = newText;
    controller.selection = TextSelection.collapsed(
      offset: start + placeholder.length,
    );

    // restore focus + keep model in sync
    focusNode?.requestFocus();
    setState(() {
      _reactionFieldValues = {
        ..._reactionFieldValues,
        target: newText,
      };
    });
  }

  Widget _outputFieldsPanel({
    required ThemeData theme,
    required AppColorPalette colors,
    required List<OutputFieldDto> outputFields,
  }) {
    if (outputFields.isEmpty) return const SizedBox.shrink();

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        _label('Variables disponibles (output)', theme, colors),
        Text(
          'Touchez une variable pour l’insérer dans le champ sélectionné (ex: {{delay}}).',
          style: theme.textTheme.bodySmall?.copyWith(color: colors.darkGrey),
        ),
        const SizedBox(height: 8),
        Wrap(
          spacing: 8,
          runSpacing: 8,
          children: outputFields.map((f) {
            final display = f.label.isNotEmpty ? f.label : f.name;
            return ActionChip(
              label: Text(display),
              onPressed: () => _insertOrCopyOutputPlaceholder(f.name),
            );
          }).toList(),
        ),
        const SizedBox(height: 16),
      ],
    );
  }

  Future<void> _createArea() async {
    if (!_canCreate) {
      setState(() {
        _createError = 'Veuillez compléter tous les champs requis.';
      });
      return;
    }

    setState(() {
      _isCreating = true;
      _createError = null;
    });

    final action = AreaActionDto(
      active: true,
      id: 0,
      provider: _actionService!.provider,
      service: _actionService!.id,
      title: _selectedAction!.title,
      type: _selectedAction!.type,
      input: _selectedAction!.fields
          .map(
            (f) => InputFieldDto(
              name: f.name,
              value: (_actionFieldValues[f.name] ?? '').trim(),
            ),
          )
          .toList(),
    );

    // Read reaction values from controllers (authoritative) if present
    final reactionInputs = _selectedReaction!.fields.map((f) {
      final controllerText = _reactionControllers[f.name]?.text;
      final v = (controllerText ?? _reactionFieldValues[f.name] ?? '').trim();
      return InputFieldDto(name: f.name, value: v);
    }).toList();

    final reaction = AreaReactionDto(
      id: 0,
      provider: _reactionService!.provider,
      service: _reactionService!.id,
      title: _selectedReaction!.title,
      input: reactionInputs,
    );

    final area = AreaDto(
      id: 0,
      name: _areaName.trim(),
      active: true,
      userId: 0,
      actions: [action],
      reactions: [reaction],
    );

    final provider = context.read<AreaProvider>();
    final res = await provider.saveArea(area);

    if (provider.error != null) {
      setState(() {
        _isCreating = false;
        _createError = provider.error;
      });
      return;
    }

    if (!mounted) return;

    if (res?.missingProviders.isNotEmpty == true) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(
            'AREA sauvegardée, mais providers manquants: ${res!.missingProviders.join(', ')}',
          ),
          backgroundColor: Colors.orange,
        ),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('AREA sauvegardée avec succès'),
          backgroundColor: Colors.green,
        ),
      );
    }

    Navigator.of(context).pop(true);
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colors = context.appColors;

    final actionServices =
        _services.where((s) => s.actions.isNotEmpty).toList(growable: false);
    final reactionServices =
        _services.where((s) => s.reactions.isNotEmpty).toList(growable: false);

    return Scaffold(
      appBar: AppBar(
        title: Text(
          _wizardStep == AreaWizardStep.action
              ? '1/3 · Action'
              : _wizardStep == AreaWizardStep.reaction
                  ? '2/3 · Réaction'
                  : '3/3 · Détails',
        ),
      ),
      body: SafeArea(
        child: _isLoading
            ? const Center(child: CircularProgressIndicator())
            : _error != null
                ? Padding(
                    padding: const EdgeInsets.all(20),
                    child: Card(
                      color: Colors.red.shade50,
                      child: Padding(
                        padding: const EdgeInsets.all(12),
                        child: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Text(_error!, style: theme.textTheme.bodySmall),
                            const SizedBox(height: 10),
                            ElevatedButton.icon(
                              onPressed: _loadFromBackend,
                              icon: const Icon(Icons.refresh),
                              label: const Text('Réessayer'),
                            ),
                          ],
                        ),
                      ),
                    ),
                  )
                : Column(
                    children: [
                      Expanded(
                        child: SingleChildScrollView(
                          padding: const EdgeInsets.all(20),
                          child: _wizardStep == AreaWizardStep.action
                              ? _buildActionStep(theme, colors, actionServices)
                              : _wizardStep == AreaWizardStep.reaction
                                  ? _buildReactionStep(
                                      theme,
                                      colors,
                                      reactionServices,
                                    )
                                  : _buildDetailsStep(theme, colors),
                        ),
                      ),
                      Padding(
                        padding: const EdgeInsets.fromLTRB(20, 0, 20, 16),
                        child: Row(
                          children: [
                            if (_wizardStep != AreaWizardStep.action)
                              Expanded(
                                child: OutlinedButton(
                                  onPressed: _prev,
                                  child: const Text('Précédent'),
                                ),
                              ),
                            if (_wizardStep != AreaWizardStep.action)
                              const SizedBox(width: 12),
                            Expanded(
                              child: ElevatedButton(
                                onPressed: _wizardStep == AreaWizardStep.details
                                    ? (_isCreating ? null : _createArea)
                                    : (_wizardStep == AreaWizardStep.action
                                        ? (_canProceedAction ? _next : null)
                                        : (_canProceedReaction ? _next : null)),
                                child: _wizardStep == AreaWizardStep.details
                                    ? (_isCreating
                                        ? const SizedBox(
                                            height: 18,
                                            width: 18,
                                            child: CircularProgressIndicator(
                                              strokeWidth: 2,
                                              color: Colors.white,
                                            ),
                                          )
                                        : const Text("Créer l'AREA"))
                                    : const Text('Suivant'),
                              ),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
      ),
    );
  }

  Widget _buildActionStep(
    ThemeData theme,
    AppColorPalette colors,
    List<_MobileService> services,
  ) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        _Section(
          title: 'Action (Déclencheur)',
          subtitle:
              'Choisissez un service, puis un déclencheur. Les services non connectés sont désactivés.',
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              _label('Service', theme, colors),
              _serviceGrid(
                services,
                selectedId: _actionService?.id,
                onTap: (s) {
                  if (!s.connected) {
                    _toastConnect();
                    return;
                  }
                  final first = s.actions.isNotEmpty ? s.actions.first : null;
                  setState(() {
                    _actionService = s;
                    _selectedAction = first;
                    _actionFieldValues =
                        first != null ? _initValues(first.fields) : {};
                  });
                },
              ),
              const SizedBox(height: 16),
              _label('Déclencheur', theme, colors),
              if (_actionService == null)
                Text(
                  'Choisissez un service.',
                  style: theme.textTheme.bodySmall
                      ?.copyWith(color: colors.darkGrey),
                )
              else
                Column(
                  children: _actionService!.actions.map((a) {
                    return RadioListTile<String>(
                      value: a.title,
                      groupValue: _selectedAction?.title,
                      title: Text(a.label),
                      dense: true,
                      contentPadding: EdgeInsets.zero,
                      onChanged: (_) {
                        setState(() {
                          _selectedAction = a;
                          _actionFieldValues = _initValues(a.fields);
                        });
                      },
                    );
                  }).toList(),
                ),
              if (_selectedAction != null) ...[
                const SizedBox(height: 16),
                _label('Paramètres', theme, colors),
                _FieldList(
                  keyPrefix: _selectedAction!.title,
                  fields: _selectedAction!.fields,
                  values: _actionFieldValues,
                  onChanged: (name, value) {
                    setState(() {
                      _actionFieldValues = {..._actionFieldValues, name: value};
                    });
                  },
                  onDatePick: _pickDateTime,
                ),
                // You can keep output fields visible here too (copy/insert will still work,
                // but insertion will target reaction fields only)
                const SizedBox(height: 8),
                _outputFieldsPanel(
                  theme: theme,
                  colors: colors,
                  outputFields: _selectedAction!.outputFields,
                ),
              ],
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildReactionStep(
    ThemeData theme,
    AppColorPalette colors,
    List<_MobileService> services,
  ) {
    return _Section(
      title: 'Réaction',
      subtitle:
          'Choisissez un service, puis une action. Les services non connectés sont désactivés.',
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _label('Service', theme, colors),
          _serviceGrid(
            services,
            selectedId: _reactionService?.id,
            onTap: (s) {
              if (!s.connected) {
                _toastConnect();
                return;
              }
              final first = s.reactions.isNotEmpty ? s.reactions.first : null;

              setState(() {
                _reactionService = s;
                _selectedReaction = first;
                _reactionFieldValues =
                    first != null ? _initValues(first.fields) : {};
              });

              if (_selectedReaction != null) {
                _initReactionEditors(_selectedReaction!.fields);
              }
            },
          ),
          const SizedBox(height: 16),
          _label('Action', theme, colors),
          if (_reactionService == null)
            Text(
              'Choisissez un service.',
              style:
                  theme.textTheme.bodySmall?.copyWith(color: colors.darkGrey),
            )
          else
            Column(
              children: _reactionService!.reactions.map((r) {
                return RadioListTile<String>(
                  value: r.title,
                  groupValue: _selectedReaction?.title,
                  title: Text(r.label),
                  dense: true,
                  contentPadding: EdgeInsets.zero,
                  onChanged: (_) {
                    setState(() {
                      _selectedReaction = r;
                      _reactionFieldValues = _initValues(r.fields);
                    });
                    _initReactionEditors(r.fields);
                  },
                );
              }).toList(),
            ),
          if (_selectedReaction != null) ...[
            const SizedBox(height: 16),
            _label('Paramètres', theme, colors),
            _FieldList(
              keyPrefix: _selectedReaction!.title,
              fields: _selectedReaction!.fields,
              values: _reactionFieldValues,
              controllers: _reactionControllers,
              focusNodes: _reactionFocusNodes,
              onChanged: (name, value) {
                setState(() {
                  _reactionFieldValues = {..._reactionFieldValues, name: value};
                });
              },
              onDatePick: _pickDateTime,
            ),
            const SizedBox(height: 8),
            if (_selectedAction != null)
              _outputFieldsPanel(
                theme: theme,
                colors: colors,
                outputFields: _selectedAction!.outputFields,
              ),
          ],
        ],
      ),
    );
  }

  Widget _buildDetailsStep(ThemeData theme, AppColorPalette colors) {
    return _Section(
      title: "Détails de l'AREA",
      subtitle: 'Ajoutez un nom et validez la création.',
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          TextField(
            onChanged: (v) => setState(() => _areaName = v),
            decoration: const InputDecoration(
              labelText: "Nom de l'AREA",
            ),
          ),
          const SizedBox(height: 12),
          if (_createError != null)
            Card(
              color: Colors.red.shade50,
              child: Padding(
                padding: const EdgeInsets.all(12),
                child: Text(
                  _createError!,
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: Colors.red.shade700,
                  ),
                ),
              ),
            ),
        ],
      ),
    );
  }

  void _toastConnect() {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Connectez ce service dans l’onglet Services.'),
        backgroundColor: Colors.orange,
      ),
    );
  }

  Widget _label(String text, ThemeData theme, AppColorPalette colors) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Text(
        text,
        style: theme.textTheme.bodySmall?.copyWith(
          color: colors.darkGrey,
          fontWeight: FontWeight.w600,
        ),
      ),
    );
  }

  Widget _serviceGrid(
    List<_MobileService> services, {
    required String? selectedId,
    required void Function(_MobileService) onTap,
  }) {
    final colors = context.appColors;

    return GridView.builder(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        childAspectRatio: 2.6,
        crossAxisSpacing: 10,
        mainAxisSpacing: 10,
      ),
      itemCount: services.length,
      itemBuilder: (context, index) {
        final s = services[index];
        final selected = selectedId == s.id;

        return InkWell(
          onTap: () => onTap(s),
          borderRadius: BorderRadius.circular(12),
          child: Opacity(
            opacity: s.connected ? 1.0 : 0.5,
            child: Container(
              height: 56,
              padding: const EdgeInsets.symmetric(horizontal: 12),
              decoration: BoxDecoration(
                color: colors.white,
                borderRadius: BorderRadius.circular(12),
                border: Border.all(
                  color: selected ? colors.deepBlue : colors.grey,
                  width: selected ? 1.4 : 1,
                ),
                boxShadow: [
                  BoxShadow(
                    color: colors.grey.withOpacity(0.2),
                    blurRadius: 10,
                    offset: const Offset(0, 4),
                  ),
                ],
              ),
              child: Center(
                child: Text(
                  s.displayName,
                  textAlign: TextAlign.center,
                  style: Theme.of(context).textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.w700,
                      ),
                ),
              ),
            ),
          ),
        );
      },
    );
  }
}

class _Section extends StatelessWidget {
  final String title;
  final String subtitle;
  final Widget child;

  const _Section({
    required this.title,
    required this.subtitle,
    required this.child,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colors = context.appColors;

    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: colors.white,
        borderRadius: BorderRadius.circular(14),
        border: Border.all(color: colors.grey),
        boxShadow: [
          BoxShadow(
            color: colors.grey.withOpacity(0.12),
            blurRadius: 12,
            offset: const Offset(0, 6),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(title, style: theme.textTheme.titleMedium),
          const SizedBox(height: 4),
          Text(
            subtitle,
            style: theme.textTheme.bodySmall?.copyWith(color: colors.darkGrey),
          ),
          const SizedBox(height: 12),
          child,
        ],
      ),
    );
  }
}

class _FieldList extends StatelessWidget {
  final String keyPrefix;
  final List<AreaFieldDefinition> fields;
  final Map<String, String> values;
  final void Function(String, String) onChanged;
  final Future<void> Function(String, void Function(String, String)) onDatePick;

  final Map<String, TextEditingController>? controllers;
  final Map<String, FocusNode>? focusNodes;

  const _FieldList({
    required this.keyPrefix,
    required this.fields,
    required this.values,
    required this.onChanged,
    required this.onDatePick,
    this.controllers,
    this.focusNodes,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colors = context.appColors;

    return Column(
      children: fields.map((field) {
        final label = field.required ? '${field.label} *' : field.label;

        final controller = controllers?[field.name];
        final focusNode = focusNodes?[field.name];
        final value = controller?.text ?? (values[field.name] ?? '');

        if (field.type == AreaFieldType.date) {
          return Padding(
            padding: const EdgeInsets.only(bottom: 12),
            child: InkWell(
              borderRadius: BorderRadius.circular(8),
              onTap: () async {
                await onDatePick(field.name, (name, iso) {
                  if (controller != null) controller.text = iso;
                  onChanged(name, iso);
                });
              },
              child: InputDecorator(
                decoration: InputDecoration(labelText: label),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      value.isEmpty ? 'Sélectionner une date' : value,
                      style: theme.textTheme.bodyMedium?.copyWith(
                        color: value.isEmpty
                            ? colors.darkGrey
                            : colors.almostBlack,
                      ),
                    ),
                    const Icon(Icons.calendar_today, size: 18),
                  ],
                ),
              ),
            ),
          );
        }

        final isNumber = field.type == AreaFieldType.number;

        return Padding(
          padding: const EdgeInsets.only(bottom: 12),
          child: TextFormField(
            key: ValueKey('$keyPrefix-${field.name}'),
            controller: controller,
            focusNode: focusNode,
            initialValue: controller == null ? value : null,
            keyboardType: isNumber ? TextInputType.number : TextInputType.text,
            inputFormatters:
                isNumber ? [FilteringTextInputFormatter.digitsOnly] : null,
            decoration: InputDecoration(labelText: label),
            onChanged: (v) => onChanged(field.name, v),
          ),
        );
      }).toList(),
    );
  }
}

class _MobileService {
  final String id;
  final String name;
  final String provider;
  final bool connected;
  final List<_MobileAction> actions;
  final List<_MobileReaction> reactions;

  _MobileService({
    required this.id,
    required this.name,
    required this.provider,
    required this.connected,
    required this.actions,
    required this.reactions,
  });

  factory _MobileService.fromConfig(ServiceConfigDto cfg,
      {required bool connected}) {
    return _MobileService(
      id: cfg.name,
      name: cfg.label.isNotEmpty ? cfg.label : cfg.name,
      provider: cfg.provider,
      connected: connected,
      actions: cfg.actions
          .map(
            (a) => _MobileAction(
              title: a.title,
              label: a.label,
              type: a.type,
              fields: a.fields.map(_fieldFromConfig).toList(),
              outputFields: a.outputFields,
            ),
          )
          .toList(),
      reactions: cfg.reactions
          .map(
            (r) => _MobileReaction(
              title: r.title,
              label: r.label,
              fields: r.fields.map(_fieldFromConfig).toList(),
            ),
          )
          .toList(),
    );
  }

  String get displayName => name;
}

class _MobileAction {
  final String title;
  final String label;
  final String type;
  final List<AreaFieldDefinition> fields;
  final List<OutputFieldDto> outputFields;

  _MobileAction({
    required this.title,
    required this.label,
    required this.type,
    required this.fields,
    required this.outputFields,
  });
}

class _MobileReaction {
  final String title;
  final String label;
  final List<AreaFieldDefinition> fields;

  _MobileReaction({
    required this.title,
    required this.label,
    required this.fields,
  });
}

AreaFieldDefinition _fieldFromConfig(ServiceFieldConfigDto f) {
  AreaFieldType type;
  switch (f.type.toLowerCase()) {
    case 'number':
      type = AreaFieldType.number;
      break;
    case 'date':
      type = AreaFieldType.date;
      break;
    default:
      type = AreaFieldType.text;
  }

  return AreaFieldDefinition(
    name: f.name,
    type: type,
    label: f.label.isNotEmpty ? f.label : f.name,
    required: f.required,
    defaultValue: f.defaultValue,
  );
}