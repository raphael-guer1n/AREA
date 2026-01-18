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

Map<String, String> _initFieldValues(List<AreaFieldDefinition> fields) {
  final out = <String, String>{};
  for (final f in fields) {
    out[f.name] = f.defaultValue ?? '';
  }
  return out;
}

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
  _MobileAction? _selectedAction;

  Map<String, String> _actionFieldValues = {};
  final List<_ReactionForm> _reactions = [];
  int _reactionId = 0;

  String _areaName = '';
  AreaWizardStep _wizardStep = AreaWizardStep.action;
  String? _createError;
  bool _isCreating = false;
  bool _showAllActionServices = true;

  @override
  void initState() {
    super.initState();
    _catalog = ServiceCatalogService();
    _connector = ServiceConnector();
    _loadFromBackend();
  }

  @override
  void dispose() {
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

  bool _reactionIsValid(_ReactionForm reaction) {
    if (reaction.service == null ||
        reaction.reaction == null ||
        !(reaction.service?.connected ?? false)) {
      return false;
    }
    return _requiredFilled(
      reaction.reaction!.fields,
      reaction.fieldValues,
    );
  }

  bool get _canProceedReaction {
    return _reactions.isNotEmpty && _reactions.every(_reactionIsValid);
  }

  bool get _canCreate {
    return _areaName.trim().isNotEmpty &&
        _canProceedAction &&
        _canProceedReaction;
  }

  void _next() {
    if (_wizardStep == AreaWizardStep.action) {
      if (_canProceedAction) {
        setState(() {
          _wizardStep = AreaWizardStep.reaction;
          _createError = null;
        });
      } else {
        setState(() {
          _createError = 'Complétez le déclencheur avant de continuer.';
        });
      }
      return;
    }
    if (_wizardStep == AreaWizardStep.reaction) {
      if (_canProceedReaction) {
        setState(() {
          _wizardStep = AreaWizardStep.details;
          _createError = null;
        });
      } else {
        setState(() {
          _createError = 'Ajoutez au moins une réaction complète.';
        });
      }
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

  void _resetActionSelection() {
    setState(() {
      _actionService = null;
      _selectedAction = null;
      _actionFieldValues = {};
      _showAllActionServices = true;
      _createError = null;
    });
  }
  Future<void> _openReactionEditor({
    required _ReactionForm reaction,
  }) async {
    final result = await Navigator.of(context).push<_ReactionForm>(
      MaterialPageRoute(
        builder: (context) => _ReactionEditorScreen(
          initial: reaction,
          services:
              _services.where((s) => s.connected && s.reactions.isNotEmpty).toList(),
          outputFields: _selectedAction?.outputFields ?? const [],
        ),
      ),
    );

    if (!mounted || result == null) return;

    setState(() {
      final index = _reactions.indexWhere((item) => item.id == result.id);
      if (index >= 0) {
        _reactions[index] = result;
      } else {
        _reactions.add(result);
      }
      _createError = null;
    });
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

    final reactions = _reactions.map((reaction) {
      final fields = reaction.reaction?.fields ?? const <AreaFieldDefinition>[];
      final inputs = fields.map((f) {
        final v = (reaction.fieldValues[f.name] ?? '').trim();
        return InputFieldDto(name: f.name, value: v);
      }).toList();
      return AreaReactionDto(
        id: 0,
        provider: reaction.service?.provider ?? '',
        service: reaction.service?.id ?? '',
        title: reaction.reaction?.title ?? '',
        input: inputs,
      );
    }).toList();

    final area = AreaDto(
      id: 0,
      name: _areaName.trim(),
      active: true,
      userId: 0,
      actions: [action],
      reactions: reactions,
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

    final actionServices = _services
        .where((s) => s.connected && s.actions.isNotEmpty)
        .toList(growable: false);
    final reactionServices = _services
        .where((s) => s.connected && s.reactions.isNotEmpty)
        .toList(growable: false);

    return Scaffold(
      appBar: AppBar(
        title: Text(
          _wizardStep == AreaWizardStep.action
              ? '1/3 · Action'
              : _wizardStep == AreaWizardStep.reaction
                  ? '2/3 · Réactions'
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
                              ? _buildReactionStep(reactionServices)
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
        if (_createError != null)
          _ErrorCard(message: _createError!),
        _Section(
          title: 'Action (Déclencheur)',
          subtitle: 'Choisissez un service puis un déclencheur.',
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const _SectionLabel(text: 'Service'),
              if (services.isEmpty)
                const _MutedPanel(text: 'Aucun service connecté disponible.')
              else if (_actionService == null || _showAllActionServices)
                _ServiceGrid(
                  services: services,
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
                          first != null ? _initFieldValues(first.fields) : {};
                      _showAllActionServices = false;
                    });
                  },
                )
              else
                _SelectedServiceCard(
                  service: _actionService!,
                  onChange: _resetActionSelection,
                ),
              const SizedBox(height: 16),
              const _SectionLabel(text: 'Déclencheur'),
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
                          _actionFieldValues = _initFieldValues(a.fields);
                        });
                      },
                    );
                  }).toList(),
                ),
              if (_selectedAction != null) ...[
                const SizedBox(height: 16),
                const _SectionLabel(text: 'Paramètres'),
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
              ],
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildReactionStep(List<_MobileService> services) {
    return _Section(
      title: 'Réactions',
      subtitle: 'Ajoutez les actions exécutées.',
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          if (_createError != null)
            Padding(
              padding: const EdgeInsets.only(bottom: 12),
              child: _ErrorCard(message: _createError!),
            ),
          if (services.isEmpty)
            _MutedPanel(
              text: 'Aucun service connecté disponible.',
            )
          else ...[
            OutlinedButton.icon(
              onPressed: () {
                final id = 'reaction-${_reactionId++}';
                _openReactionEditor(
                  reaction: _ReactionForm(id: id),
                );
              },
              icon: const Icon(Icons.add),
              label: const Text('Ajouter une réaction'),
            ),
            const SizedBox(height: 12),
            if (_reactions.isEmpty)
              _MutedPanel(
                text: 'Aucune réaction ajoutée.',
              )
            else
              Column(
                children: _reactions.asMap().entries.map((entry) {
                  final index = entry.key;
                  final reaction = entry.value;
                  final label = reaction.reaction?.label ??
                      reaction.service?.displayName ??
                      'Réaction ${index + 1}';
                  final isValid = _reactionIsValid(reaction);

                  return _ReactionSummaryCard(
                    key: ValueKey(reaction.id),
                    label: label,
                    service: reaction.service,
                    isValid: isValid,
                    onTap: () => _openReactionEditor(reaction: reaction),
                    onDelete: () {
                      setState(() {
                        _removeReaction(reaction);
                      });
                    },
                  );
                }).toList(),
              ),
          ],
        ],
      ),
    );
  }

  Widget _buildDetailsStep(ThemeData theme, AppColorPalette colors) {
    final actionSummary = _actionService == null || _selectedAction == null
        ? 'Action non définie'
        : '${_actionService!.displayName} · ${_selectedAction!.label}';
    final actionItems = _selectedAction == null
        ? const <_SummaryItem>[]
        : _buildFieldSummaryItems(_selectedAction!.fields, _actionFieldValues);

    return _Section(
      title: "Détails de l'AREA",
      subtitle: 'Nom et validation.',
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
          const SizedBox(height: 12),
          _SummaryCard(
            title: 'Résumé',
            children: [
              _SummaryGroup(
                title: 'Action',
                subtitle: actionSummary,
                items: actionItems,
              ),
              const SizedBox(height: 12),
              if (_reactions.isEmpty)
                Text(
                  'Aucune réaction ajoutée.',
                  style: theme.textTheme.bodySmall
                      ?.copyWith(color: colors.darkGrey),
                )
              else
                Column(
                  children: _reactions.asMap().entries.map((entry) {
                    final index = entry.key + 1;
                    final reaction = entry.value;
                    final subtitle =
                        reaction.service == null || reaction.reaction == null
                            ? 'Réaction non définie'
                            : '${reaction.service!.displayName} · ${reaction.reaction!.label}';
                    final items = reaction.reaction == null
                        ? const <_SummaryItem>[]
                        : _buildFieldSummaryItems(
                            reaction.reaction!.fields,
                            reaction.fieldValues,
                          );

                    return Padding(
                      padding: const EdgeInsets.only(bottom: 12),
                      child: _SummaryGroup(
                        title: 'Réaction $index',
                        subtitle: subtitle,
                        items: items,
                      ),
                    );
                  }).toList(),
                ),
            ],
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

  void _removeReaction(_ReactionForm reaction) {
    _reactions.removeWhere((item) => item.id == reaction.id);
    _createError = null;
  }

}

List<_SummaryItem> _buildFieldSummaryItems(
  List<AreaFieldDefinition> fields,
  Map<String, String> values,
) {
  return fields.map((field) {
    final value = values[field.name]?.trim() ?? '';
    return _SummaryItem(
      label: field.label,
      value: value.isEmpty ? '—' : value,
    );
  }).toList();
}

class _SelectedServiceCard extends StatelessWidget {
  final _MobileService service;
  final VoidCallback onChange;

  const _SelectedServiceCard({
    required this.service,
    required this.onChange,
  });

  @override
  Widget build(BuildContext context) {
    final colors = context.appColors;

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 12),
      decoration: BoxDecoration(
        color: colors.white,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: colors.grey),
        boxShadow: [
          BoxShadow(
            color: colors.grey.withOpacity(0.18),
            blurRadius: 10,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Row(
        children: [
          _ServiceLogo(service: service, size: 34),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              service.displayName,
              style: Theme.of(context).textTheme.titleSmall?.copyWith(
                    fontWeight: FontWeight.w700,
                  ),
            ),
          ),
          TextButton(
            onPressed: onChange,
            child: const Text('Changer'),
          ),
        ],
      ),
    );
  }
}

class _ServiceGrid extends StatelessWidget {
  final List<_MobileService> services;
  final String? selectedId;
  final void Function(_MobileService) onTap;

  const _ServiceGrid({
    required this.services,
    required this.selectedId,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
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
                  color: selected ? colors.midBlue : colors.grey,
                  width: selected ? 1.6 : 1,
                ),
                boxShadow: [
                  BoxShadow(
                    color: colors.grey.withOpacity(0.18),
                    blurRadius: 10,
                    offset: const Offset(0, 4),
                  ),
                ],
              ),
              child: Row(
                children: [
                  _ServiceLogo(service: s, size: 30),
                  const SizedBox(width: 10),
                  Expanded(
                    child: Text(
                      s.displayName,
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                      style: Theme.of(context).textTheme.titleSmall?.copyWith(
                            fontWeight: FontWeight.w700,
                            color: colors.almostBlack,
                          ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        );
      },
    );
  }
}

class _ServiceLogo extends StatelessWidget {
  final _MobileService service;
  final double size;

  const _ServiceLogo({required this.service, required this.size});

  @override
  Widget build(BuildContext context) {
    final colors = context.appColors;
    final initials = service.displayName.isNotEmpty
        ? service.displayName.trim().substring(0, 1).toUpperCase()
        : '?';

    if (service.logoUrl.isEmpty) {
      return CircleAvatar(
        radius: size / 2,
        backgroundColor: colors.deepBlue.withOpacity(0.1),
        child: Text(
          initials,
          style: Theme.of(context).textTheme.labelLarge?.copyWith(
                color: colors.deepBlue,
                fontWeight: FontWeight.w700,
              ),
        ),
      );
    }

    return ClipRRect(
      borderRadius: BorderRadius.circular(size / 2),
      child: Image.network(
        service.logoUrl,
        width: size,
        height: size,
        fit: BoxFit.cover,
        errorBuilder: (_, __, ___) => CircleAvatar(
          radius: size / 2,
          backgroundColor: colors.deepBlue.withOpacity(0.1),
          child: Text(
            initials,
            style: Theme.of(context).textTheme.labelLarge?.copyWith(
                  color: colors.deepBlue,
                  fontWeight: FontWeight.w700,
                ),
          ),
        ),
      ),
    );
  }
}

class _ReactionSummaryCard extends StatelessWidget {
  final String label;
  final _MobileService? service;
  final bool isValid;
  final VoidCallback onTap;
  final VoidCallback onDelete;

  const _ReactionSummaryCard({
    super.key,
    required this.label,
    required this.service,
    required this.isValid,
    required this.onTap,
    required this.onDelete,
  });

  @override
  Widget build(BuildContext context) {
    final colors = context.appColors;
    final badgeColor = isValid ? Colors.green : Colors.orange;
    final badgeText = isValid ? 'Complet' : 'Incomplet';

    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(14),
        child: Container(
          padding: const EdgeInsets.all(14),
          decoration: BoxDecoration(
            color: colors.white,
            borderRadius: BorderRadius.circular(14),
            border: Border.all(color: colors.grey),
            boxShadow: [
              BoxShadow(
                color: colors.grey.withOpacity(0.16),
                blurRadius: 10,
                offset: const Offset(0, 4),
              ),
            ],
          ),
          child: Row(
            children: [
              if (service != null) _ServiceLogo(service: service!, size: 34),
              if (service != null) const SizedBox(width: 10),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      label,
                      style: Theme.of(context).textTheme.titleSmall?.copyWith(
                            fontWeight: FontWeight.w700,
                          ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      service?.displayName ?? 'Service non défini',
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                            color: colors.darkGrey,
                          ),
                    ),
                  ],
                ),
              ),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
                decoration: BoxDecoration(
                  color: badgeColor.withOpacity(0.16),
                  borderRadius: BorderRadius.circular(20),
                  border: Border.all(color: badgeColor.withOpacity(0.5)),
                ),
                child: Text(
                  badgeText,
                  style: Theme.of(context).textTheme.labelSmall?.copyWith(
                        color: badgeColor,
                        fontWeight: FontWeight.w700,
                      ),
                ),
              ),
              IconButton(
                onPressed: onDelete,
                icon: const Icon(Icons.delete_outline),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _MutedPanel extends StatelessWidget {
  final String text;

  const _MutedPanel({required this.text});

  @override
  Widget build(BuildContext context) {
    final colors = context.appColors;

    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: colors.white,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: colors.grey),
      ),
      child: Text(
        text,
        style: Theme.of(context).textTheme.bodySmall?.copyWith(
              color: colors.darkGrey,
            ),
      ),
    );
  }
}

class _SectionLabel extends StatelessWidget {
  final String text;

  const _SectionLabel({required this.text});

  @override
  Widget build(BuildContext context) {
    final colors = context.appColors;

    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Text(
        text,
        style: Theme.of(context).textTheme.bodySmall?.copyWith(
              color: colors.darkGrey,
              fontWeight: FontWeight.w600,
            ),
      ),
    );
  }
}

class _ReactionEditorScreen extends StatefulWidget {
  final _ReactionForm initial;
  final List<_MobileService> services;
  final List<OutputFieldDto> outputFields;

  const _ReactionEditorScreen({
    required this.initial,
    required this.services,
    required this.outputFields,
  });

  @override
  State<_ReactionEditorScreen> createState() => _ReactionEditorScreenState();
}

class _ReactionEditorScreenState extends State<_ReactionEditorScreen> {
  _MobileService? _service;
  _MobileReaction? _reaction;
  Map<String, String> _values = {};
  final Map<String, TextEditingController> _controllers = {};
  final Map<String, FocusNode> _focusNodes = {};
  String? _focusedField;
  String? _error;

  @override
  void initState() {
    super.initState();
    _service = widget.initial.service;
    _reaction = widget.initial.reaction;
    _values = Map<String, String>.from(widget.initial.fieldValues);
    if (_reaction != null) {
      _initControllers(_reaction!.fields);
    }
  }

  @override
  void dispose() {
    for (final c in _controllers.values) {
      c.dispose();
    }
    for (final f in _focusNodes.values) {
      f.dispose();
    }
    super.dispose();
  }

  void _initControllers(List<AreaFieldDefinition> fields) {
    for (final c in _controllers.values) {
      c.dispose();
    }
    for (final f in _focusNodes.values) {
      f.dispose();
    }
    _controllers.clear();
    _focusNodes.clear();
    _focusedField = null;

    for (final field in fields) {
      final initialValue = _values[field.name] ?? (field.defaultValue ?? '');
      final controller = TextEditingController(text: initialValue);
      final focusNode = FocusNode();

      focusNode.addListener(() {
        if (focusNode.hasFocus) {
          setState(() {
            _focusedField = field.name;
          });
        }
      });

      _controllers[field.name] = controller;
      _focusNodes[field.name] = focusNode;
    }
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

  Future<void> _insertToken(String outputName) async {
    final placeholder = '{{${outputName.trim()}}}';
    final target = _focusedField;

    if (target == null || !_controllers.containsKey(target)) {
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

    final controller = _controllers[target]!;
    final focusNode = _focusNodes[target];

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

    focusNode?.requestFocus();
    setState(() {
      _values = {
        ..._values,
        target: newText,
      };
    });
  }

  void _save() {
    if (_service == null) {
      setState(() {
        _error = 'Choisissez un service.';
      });
      return;
    }
    if (_reaction == null) {
      setState(() {
        _error = 'Choisissez une action.';
      });
      return;
    }

    if (!_requiredFilled(_reaction!.fields, _values)) {
      setState(() {
        _error = 'Complétez les champs obligatoires.';
      });
      return;
    }

    final updated = _ReactionForm(
      id: widget.initial.id,
      service: _service,
      reaction: _reaction,
      fieldValues: Map<String, String>.from(_values),
    );

    Navigator.of(context).pop(updated);
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colors = context.appColors;
    final title = widget.initial.reaction == null
        ? 'Nouvelle réaction'
        : 'Modifier réaction';

    return Scaffold(
      appBar: AppBar(title: Text(title)),
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(20),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              if (_error != null)
                Padding(
                  padding: const EdgeInsets.only(bottom: 12),
                  child: _ErrorCard(message: _error!),
                ),
              const _SectionLabel(text: 'Service'),
              if (widget.services.isEmpty)
                _MutedPanel(text: 'Aucun service connecté disponible.')
              else if (_service == null)
                _ServiceGrid(
                  services: widget.services,
                  selectedId: _service?.id,
                  onTap: (s) {
                    final first =
                        s.reactions.isNotEmpty ? s.reactions.first : null;
                    setState(() {
                      _service = s;
                      _reaction = first;
                      _values =
                          first != null ? _initFieldValues(first.fields) : {};
                      _error = null;
                    });
                    if (_reaction != null) {
                      _initControllers(_reaction!.fields);
                    }
                  },
                )
              else
                _SelectedServiceCard(
                  service: _service!,
                  onChange: () {
                    setState(() {
                      _service = null;
                      _reaction = null;
                      _values = {};
                      _controllers.clear();
                      _focusNodes.clear();
                    });
                  },
                ),
              const SizedBox(height: 16),
              const _SectionLabel(text: 'Action'),
              if (_service == null)
                Text(
                  'Choisissez un service.',
                  style:
                      theme.textTheme.bodySmall?.copyWith(color: colors.darkGrey),
                )
              else
                Column(
                  children: _service!.reactions.map((r) {
                    return RadioListTile<String>(
                      value: r.title,
                      groupValue: _reaction?.title,
                      title: Text(r.label),
                      dense: true,
                      contentPadding: EdgeInsets.zero,
                      onChanged: (_) {
                        setState(() {
                          _reaction = r;
                          _values = _initFieldValues(r.fields);
                          _error = null;
                        });
                        _initControllers(r.fields);
                      },
                    );
                  }).toList(),
                ),
              if (_reaction != null) ...[
                const SizedBox(height: 16),
                const _SectionLabel(text: 'Paramètres'),
                _FieldList(
                  keyPrefix: '${widget.initial.id}-${_reaction!.title}',
                  fields: _reaction!.fields,
                  values: _values,
                  controllers: _controllers,
                  focusNodes: _focusNodes,
                  onChanged: (name, value) {
                    setState(() {
                      _values = {..._values, name: value};
                    });
                  },
                  onDatePick: _pickDateTime,
                ),
                const SizedBox(height: 12),
                _OutputFieldsPanel(
                  outputFields: widget.outputFields,
                  onTapToken: _insertToken,
                ),
              ],
            ],
          ),
        ),
      ),
      bottomNavigationBar: Padding(
        padding: const EdgeInsets.fromLTRB(20, 0, 20, 16),
        child: ElevatedButton(
          onPressed: (_service != null &&
                  _reaction != null &&
                  _requiredFilled(_reaction!.fields, _values))
              ? _save
              : null,
          child: const Text('Suivant'),
        ),
      ),
    );
  }
}

class _OutputFieldsPanel extends StatelessWidget {
  final List<OutputFieldDto> outputFields;
  final void Function(String) onTapToken;

  const _OutputFieldsPanel({
    required this.outputFields,
    required this.onTapToken,
  });

  @override
  Widget build(BuildContext context) {
    if (outputFields.isEmpty) return const SizedBox.shrink();

    final colors = context.appColors;

    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: colors.white,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: colors.grey),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Variables',
            style: Theme.of(context).textTheme.titleSmall?.copyWith(
                  fontWeight: FontWeight.w700,
                ),
          ),
          const SizedBox(height: 6),
          Text(
            'Touchez pour insérer',
            style: Theme.of(context).textTheme.bodySmall?.copyWith(
                  color: colors.darkGrey,
                ),
          ),
          const SizedBox(height: 10),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: outputFields.map((field) {
              final label = field.label.isNotEmpty ? field.label : field.name;
              return InkWell(
                onTap: () => onTapToken(field.name),
                borderRadius: BorderRadius.circular(18),
                child: Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                  decoration: BoxDecoration(
                    color: colors.midBlue,
                    borderRadius: BorderRadius.circular(18),
                  ),
                  child: Text(
                    label,
                    style: Theme.of(context).textTheme.labelSmall?.copyWith(
                          color: Colors.white,
                          fontWeight: FontWeight.w700,
                        ),
                  ),
                ),
              );
            }).toList(),
          ),
        ],
      ),
    );
  }
}

class _SummaryCard extends StatelessWidget {
  final String title;
  final List<Widget> children;

  const _SummaryCard({
    required this.title,
    required this.children,
  });

  @override
  Widget build(BuildContext context) {
    final colors = context.appColors;

    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: colors.white,
        borderRadius: BorderRadius.circular(14),
        border: Border.all(color: colors.midBlue.withOpacity(0.4)),
        boxShadow: [
          BoxShadow(
            color: colors.grey.withOpacity(0.18),
            blurRadius: 10,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                width: 4,
                height: 18,
                decoration: BoxDecoration(
                  color: colors.midBlue,
                  borderRadius: BorderRadius.circular(8),
                ),
              ),
              const SizedBox(width: 8),
              Text(
                title,
                style: Theme.of(context).textTheme.titleSmall?.copyWith(
                      fontWeight: FontWeight.w700,
                      color: colors.deepBlue,
                    ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          ...children,
        ],
      ),
    );
  }
}

class _SummaryGroup extends StatelessWidget {
  final String title;
  final String subtitle;
  final List<_SummaryItem> items;

  const _SummaryGroup({
    required this.title,
    required this.subtitle,
    required this.items,
  });

  @override
  Widget build(BuildContext context) {
    final colors = context.appColors;
    final theme = Theme.of(context);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          title,
          style: theme.textTheme.bodySmall?.copyWith(
            color: colors.darkGrey,
            fontWeight: FontWeight.w600,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          subtitle,
          style: theme.textTheme.bodyMedium?.copyWith(
            color: colors.almostBlack,
            fontWeight: FontWeight.w600,
          ),
        ),
        const SizedBox(height: 6),
        if (items.isEmpty)
          Text(
            'Aucun paramètre',
            style: theme.textTheme.bodySmall?.copyWith(
              color: colors.darkGrey,
            ),
          )
        else
          Column(
            children: items.map((item) {
              return Padding(
                padding: const EdgeInsets.only(bottom: 4),
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Container(
                      width: 6,
                      height: 6,
                      margin: const EdgeInsets.only(top: 6, right: 8),
                      decoration: BoxDecoration(
                        color: colors.midBlue,
                        shape: BoxShape.circle,
                      ),
                    ),
                    Expanded(
                      child: Text(
                        '${item.label} : ${item.value}',
                        style: theme.textTheme.bodySmall?.copyWith(
                          color: colors.almostBlack,
                        ),
                      ),
                    ),
                  ],
                ),
              );
            }).toList(),
          ),
      ],
    );
  }
}

class _SummaryItem {
  final String label;
  final String value;

  const _SummaryItem({
    required this.label,
    required this.value,
  });
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
        color: colors.lightGrey,
        borderRadius: BorderRadius.circular(14),
        border: Border.all(color: colors.grey),
        boxShadow: [
          BoxShadow(
            color: colors.grey.withOpacity(0.18),
            blurRadius: 12,
            offset: const Offset(0, 6),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                width: 4,
                height: 18,
                decoration: BoxDecoration(
                  color: colors.midBlue,
                  borderRadius: BorderRadius.circular(8),
                ),
              ),
              const SizedBox(width: 8),
              Text(
                title,
                style: theme.textTheme.titleMedium?.copyWith(
                      color: colors.deepBlue,
                      fontWeight: FontWeight.w700,
                    ),
              ),
            ],
          ),
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

class _ErrorCard extends StatelessWidget {
  final String message;

  const _ErrorCard({required this.message});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Card(
      color: Colors.red.shade50,
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Text(
          message,
          style: theme.textTheme.bodySmall
              ?.copyWith(color: Colors.red.shade700),
        ),
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
  final String logoUrl;
  final bool connected;
  final List<_MobileAction> actions;
  final List<_MobileReaction> reactions;

  _MobileService({
    required this.id,
    required this.name,
    required this.provider,
    required this.logoUrl,
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
      logoUrl: cfg.logoUrl.isNotEmpty
          ? cfg.logoUrl
          : (cfg.iconUrl.isNotEmpty ? cfg.iconUrl : ''),
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

class _ReactionForm {
  final String id;
  _MobileService? service;
  _MobileReaction? reaction;
  Map<String, String> fieldValues;

  _ReactionForm({
    required this.id,
    this.service,
    this.reaction,
    Map<String, String>? fieldValues,
  }) : fieldValues = fieldValues ?? {};
}
