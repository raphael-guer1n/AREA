import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:provider/provider.dart';
import '../../data/area_services.dart';
import '../../models/area_definitions.dart';
import '../../models/area_model.dart';
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
  late final ServiceCatalogService _catalogService;
  late final ServiceConnector _serviceConnector;

  List<AreaServiceDefinition> _services = [];
  bool _isLoadingServices = true;
  String? _servicesError;

  AreaServiceDefinition? _actionService;
  AreaServiceDefinition? _reactionService;
  AreaActionDefinition? _selectedAction;
  AreaReactionDefinition? _selectedReaction;
  Map<String, String> _actionFieldValues = {};
  Map<String, String> _reactionFieldValues = {};
  String _areaName = '';
  AreaWizardStep _wizardStep = AreaWizardStep.action;
  String? _createError;
  bool _isCreating = false;

  @override
  void initState() {
    super.initState();
    _catalogService = ServiceCatalogService();
    _serviceConnector = ServiceConnector();
    _loadServices();
  }

  Future<void> _loadServices() async {
    setState(() {
      _isLoadingServices = true;
      _servicesError = null;
    });

    try {
      final serviceIds = await _catalogService.fetchServices();
      final uniqueServiceIds = <String>{...serviceIds, 'timer'}.toList();

      final authProvider = context.read<AuthProvider>();
      final rawUserId = authProvider.user?['id'];
      int? userId;
      if (rawUserId is int) {
        userId = rawUserId;
      } else if (rawUserId != null) {
        userId = int.tryParse(rawUserId.toString());
      }
      final Map<String, bool> statusByService = {};

      if (userId != null) {
        final statuses = await _serviceConnector.fetchServices(userId);
        for (final status in statuses) {
          statusByService[status.name] = status.isConnected;
        }
      }

      final mappedServices = uniqueServiceIds.map((serviceId) {
        final base = areaServiceCatalog.firstWhere(
          (service) => service.id == serviceId,
          orElse: () => AreaServiceDefinition(
            id: serviceId,
            name: _formatServiceName(serviceId),
          ),
        );
        final isConnected =
            statusByService[serviceId] ?? (serviceId == 'timer');
        return base.copyWith(connected: isConnected);
      }).toList();

      setState(() {
        _services = mappedServices;
      });
    } catch (e) {
      setState(() {
        _servicesError = e.toString();
        _services = [];
      });
    } finally {
      setState(() {
        _isLoadingServices = false;
      });
    }
  }

  String _formatServiceName(String serviceId) {
    if (serviceId.isEmpty) return serviceId;
    return serviceId[0].toUpperCase() + serviceId.substring(1);
  }

  Map<String, String> _initializeFieldValues(
      List<AreaFieldDefinition> fields) {
    final values = <String, String>{};
    for (final field in fields) {
      values[field.name] = field.defaultValue ?? '';
    }
    return values;
  }

  bool _areRequiredFieldsFilled(
    List<AreaFieldDefinition> fields,
    Map<String, String> values,
  ) {
    return fields.every((field) {
      if (!field.required) return true;
      final value = values[field.name];
      return value != null && value.trim().isNotEmpty;
    });
  }

  bool get _canProceedAction {
    return _actionService != null &&
        _actionService!.connected &&
        _selectedAction != null &&
        _areRequiredFieldsFilled(
          _selectedAction?.fields ?? [],
          _actionFieldValues,
        );
  }

  bool get _canProceedReaction {
    return _reactionService != null &&
        _reactionService!.connected &&
        _selectedReaction != null &&
        _areRequiredFieldsFilled(
          _selectedReaction?.fields ?? [],
          _reactionFieldValues,
        );
  }

  bool get _canCreate {
    return _actionService != null &&
        _reactionService != null &&
        _selectedAction != null &&
        _selectedReaction != null &&
        _areaName.trim().isNotEmpty &&
        _areRequiredFieldsFilled(
          _selectedAction?.fields ?? [],
          _actionFieldValues,
        ) &&
        _areRequiredFieldsFilled(
          _selectedReaction?.fields ?? [],
          _reactionFieldValues,
        );
  }

  void _goToNextStep() {
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

  void _goToPreviousStep() {
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

  Future<void> _handleCreateArea() async {
    if (!_canCreate) {
      setState(() {
        _createError = 'Veuillez compléter tous les champs requis.';
      });
      return;
    }

    final delayValue =
        int.tryParse(_actionFieldValues['delay'] ?? '0') ?? 0;
    final startTimeValue = _reactionFieldValues['start_time'] ?? '';
    final endTimeValue = _reactionFieldValues['end_time'] ?? '';
    final summaryValue = (_reactionFieldValues['summary'] ?? '').trim();
    final descriptionValue =
        (_reactionFieldValues['description'] ?? '').trim();
    final fallbackSummary = summaryValue.isEmpty ? _areaName.trim() : summaryValue;

    DateTime startTime;
    DateTime endTime;
    try {
      startTime = DateTime.parse(startTimeValue);
      endTime = DateTime.parse(endTimeValue);
    } catch (_) {
      setState(() {
        _createError = 'Les dates de début et fin sont invalides.';
      });
      return;
    }

    setState(() {
      _isCreating = true;
      _createError = null;
    });

    final provider = context.read<AreaProvider>();
    final req = CreateEventRequest(
      delay: delayValue,
      event: EventModel(
        startTime: startTime.toUtc().toIso8601String(),
        endTime: endTime.toUtc().toIso8601String(),
        summary: fallbackSummary,
        description: descriptionValue,
      ),
    );

    await provider.createEvent(req);

    if (provider.error != null) {
      setState(() {
        _isCreating = false;
        _createError = provider.error;
      });
      return;
    }

    final created = CreatedArea(
      id: 'area-${DateTime.now().millisecondsSinceEpoch}',
      name: _areaName.trim(),
      summary: fallbackSummary,
      startTime: startTime.toIso8601String(),
      endTime: endTime.toIso8601String(),
      delay: delayValue,
      actionService: _actionService?.name ?? '',
      reactionService: _reactionService?.name ?? '',
      actionName: _selectedAction?.label ?? '',
      reactionName: _selectedReaction?.label ?? '',
      gradient: pickRandomGradient(),
    );

    if (!mounted) return;
    Navigator.of(context).pop(created);
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
    onChanged(fieldName, dt.toIso8601String());
  }

  String _formatDateTime(DateTime dt) {
    String two(int value) => value.toString().padLeft(2, '0');
    return '${dt.year}-${two(dt.month)}-${two(dt.day)} ${two(dt.hour)}:${two(dt.minute)}';
  }

  String _formatFieldValue(
    AreaFieldDefinition field,
    Map<String, String> values,
  ) {
    final value = values[field.name] ?? '';
    if (value.trim().isEmpty) return '—';
    if (field.type == AreaFieldType.date) {
      try {
        return _formatDateTime(DateTime.parse(value));
      } catch (_) {
        return value;
      }
    }
    return value;
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final connectedServices =
        _services.where((service) => service.connected).toList();
    final steps = const [
      _WizardStepInfo(
        step: AreaWizardStep.action,
        title: 'Action',
        description: 'Déclencheur',
      ),
      _WizardStepInfo(
        step: AreaWizardStep.reaction,
        title: 'Réaction',
        description: 'Action exécutée',
      ),
      _WizardStepInfo(
        step: AreaWizardStep.details,
        title: 'Détails',
        description: 'Planification',
      ),
    ];
    final currentIndex =
        steps.indexWhere((step) => step.step == _wizardStep);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Créer une AREA'),
      ),
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.all(24),
          children: [
            if (_isLoadingServices)
              const LinearProgressIndicator(minHeight: 2),
            const SizedBox(height: 12),
            Text(
              'Composez votre automation',
              style: theme.textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(height: 6),
            Text(
              'Sélectionnez un déclencheur et une réaction pour structurer votre AREA.',
              style: theme.textTheme.bodySmall?.copyWith(
                color: AppColors.darkGrey,
              ),
            ),
            const SizedBox(height: 16),
            Row(
              children: steps.asMap().entries.map((entry) {
                final index = entry.key;
                final step = entry.value;
                final isActive = step.step == _wizardStep;
                final isDone = index < currentIndex;
                return Expanded(
                  child: Container(
                    margin: EdgeInsets.only(
                      right: index == steps.length - 1 ? 0 : 8,
                    ),
                    padding: const EdgeInsets.symmetric(vertical: 10),
                    decoration: BoxDecoration(
                      color: isActive
                          ? AppColors.deepBlue.withOpacity(0.08)
                          : AppColors.white,
                      borderRadius: BorderRadius.circular(12),
                      border: Border.all(
                        color: isActive || isDone
                            ? AppColors.deepBlue
                            : AppColors.grey,
                      ),
                    ),
                    child: Column(
                      children: [
                        Text(
                          '${index + 1}',
                          style: theme.textTheme.bodyMedium?.copyWith(
                            fontWeight: FontWeight.w700,
                            color: isActive || isDone
                                ? AppColors.deepBlue
                                : AppColors.darkGrey,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          step.title,
                          style: theme.textTheme.bodySmall?.copyWith(
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                        Text(
                          step.description,
                          style: theme.textTheme.labelSmall?.copyWith(
                            color: AppColors.darkGrey,
                          ),
                        ),
                      ],
                    ),
                  ),
                );
              }).toList(),
            ),
            const SizedBox(height: 16),
            if (_servicesError != null)
              Card(
                color: Colors.red.shade50,
                child: Padding(
                  padding: const EdgeInsets.all(12),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        _servicesError!,
                        style: theme.textTheme.bodySmall?.copyWith(
                          color: Colors.red.shade700,
                        ),
                      ),
                      const SizedBox(height: 8),
                      TextButton.icon(
                        onPressed: _loadServices,
                        icon: const Icon(Icons.refresh),
                        label: const Text('Réessayer'),
                      ),
                    ],
                  ),
                ),
              ),
            if (_wizardStep == AreaWizardStep.action)
              _AreaSection(
                title: 'Action (Déclencheur)',
                subtitle:
                    'Choisissez un service connecté puis le déclencheur.',
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _SectionLabel(label: 'Service'),
                    if (connectedServices.isEmpty)
                      Text(
                        'Aucun service connecté. Connectez-en un depuis l\'onglet Services.',
                        style: theme.textTheme.bodySmall?.copyWith(
                          color: AppColors.darkGrey,
                        ),
                      )
                    else
                      Wrap(
                        spacing: 8,
                        runSpacing: 8,
                        children: connectedServices.map((service) {
                          final isSelected = _actionService?.id == service.id;
                          return ChoiceChip(
                            label: Text(service.displayName),
                            selected: isSelected,
                            onSelected: (_) {
                              final defaultAction =
                                  service.actions.isNotEmpty
                                      ? service.actions.first
                                      : null;
                              setState(() {
                                _actionService = service;
                                _selectedAction = defaultAction;
                                _actionFieldValues = defaultAction != null
                                    ? _initializeFieldValues(
                                        defaultAction.fields,
                                      )
                                    : {};
                              });
                            },
                          );
                        }).toList(),
                      ),
                    const SizedBox(height: 16),
                    _SectionLabel(label: 'Déclencheur'),
                    if (_actionService == null)
                      Text(
                        'Choisissez un service pour voir ses déclencheurs.',
                        style: theme.textTheme.bodySmall?.copyWith(
                          color: AppColors.darkGrey,
                        ),
                      )
                    else if (_actionService!.actions.isEmpty)
                      Text(
                        'Aucun déclencheur disponible pour ce service.',
                        style: theme.textTheme.bodySmall?.copyWith(
                          color: AppColors.darkGrey,
                        ),
                      )
                    else
                      Column(
                        children: _actionService!.actions.map((action) {
                          return RadioListTile<String>(
                            value: action.id,
                            groupValue: _selectedAction?.id,
                            title: Text(action.label),
                            dense: true,
                            contentPadding: EdgeInsets.zero,
                            onChanged: (_) {
                              setState(() {
                                _selectedAction = action;
                                _actionFieldValues =
                                    _initializeFieldValues(action.fields);
                              });
                            },
                          );
                        }).toList(),
                      ),
                    if (_selectedAction != null) ...[
                      const SizedBox(height: 16),
                      _SectionLabel(label: 'Paramètres du déclencheur'),
                      _FieldList(
                        keyPrefix: _selectedAction!.id,
                        fields: _selectedAction!.fields,
                        values: _actionFieldValues,
                        onChanged: (name, value) {
                          setState(() {
                            _actionFieldValues = {
                              ..._actionFieldValues,
                              name: value,
                            };
                          });
                        },
                        onDatePick: _pickDateTime,
                        formatDate: _formatDateTime,
                      ),
                    ],
                  ],
                ),
              ),
            if (_wizardStep == AreaWizardStep.reaction)
              _AreaSection(
                title: 'Réaction',
                subtitle:
                    'Sélectionnez le service qui exécutera l\'action après le déclencheur.',
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _SectionLabel(label: 'Service'),
                    if (connectedServices.isEmpty)
                      Text(
                        'Aucun service connecté. Connectez-en un depuis l\'onglet Services.',
                        style: theme.textTheme.bodySmall?.copyWith(
                          color: AppColors.darkGrey,
                        ),
                      )
                    else
                      Wrap(
                        spacing: 8,
                        runSpacing: 8,
                        children: connectedServices.map((service) {
                          final isSelected = _reactionService?.id == service.id;
                          return ChoiceChip(
                            label: Text(service.displayName),
                            selected: isSelected,
                            onSelected: (_) {
                              final defaultReaction =
                                  service.reactions.isNotEmpty
                                      ? service.reactions.first
                                      : null;
                              setState(() {
                                _reactionService = service;
                                _selectedReaction = defaultReaction;
                                _reactionFieldValues = defaultReaction != null
                                    ? _initializeFieldValues(
                                        defaultReaction.fields,
                                      )
                                    : {};
                              });
                            },
                          );
                        }).toList(),
                      ),
                    const SizedBox(height: 16),
                    _SectionLabel(label: 'Action'),
                    if (_reactionService == null)
                      Text(
                        'Choisissez un service pour voir ses actions.',
                        style: theme.textTheme.bodySmall?.copyWith(
                          color: AppColors.darkGrey,
                        ),
                      )
                    else if (_reactionService!.reactions.isEmpty)
                      Text(
                        'Aucune action disponible pour ce service.',
                        style: theme.textTheme.bodySmall?.copyWith(
                          color: AppColors.darkGrey,
                        ),
                      )
                    else
                      Column(
                        children: _reactionService!.reactions.map((reaction) {
                          return RadioListTile<String>(
                            value: reaction.id,
                            groupValue: _selectedReaction?.id,
                            title: Text(reaction.label),
                            dense: true,
                            contentPadding: EdgeInsets.zero,
                            onChanged: (_) {
                              setState(() {
                                _selectedReaction = reaction;
                                _reactionFieldValues =
                                    _initializeFieldValues(reaction.fields);
                              });
                            },
                          );
                        }).toList(),
                      ),
                    if (_selectedReaction != null) ...[
                      const SizedBox(height: 16),
                      _SectionLabel(label: 'Paramètres de la réaction'),
                      _FieldList(
                        keyPrefix: _selectedReaction!.id,
                        fields: _selectedReaction!.fields,
                        values: _reactionFieldValues,
                        onChanged: (name, value) {
                          setState(() {
                            _reactionFieldValues = {
                              ..._reactionFieldValues,
                              name: value,
                            };
                          });
                        },
                        onDatePick: _pickDateTime,
                        formatDate: _formatDateTime,
                      ),
                    ],
                  ],
                ),
              ),
            if (_wizardStep == AreaWizardStep.details)
              Column(
                children: [
                  _AreaSection(
                    title: 'Détails de l\'AREA',
                    subtitle:
                        'Ajoutez un nom et vérifiez les paramètres.',
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        TextField(
                          onChanged: (value) {
                            setState(() {
                              _areaName = value;
                            });
                          },
                          decoration: const InputDecoration(
                            labelText: 'Nom de l\'AREA',
                            hintText: 'Démo marketing',
                          ),
                        ),
                        const SizedBox(height: 12),
                        Container(
                          padding: const EdgeInsets.all(12),
                          decoration: BoxDecoration(
                            color: AppColors.lightGrey,
                            borderRadius: BorderRadius.circular(8),
                            border: Border.all(color: AppColors.grey),
                          ),
                          child: Text(
                            'Vérifiez les paramètres du déclencheur et de la réaction avant de créer l\'AREA.',
                            style: theme.textTheme.bodySmall?.copyWith(
                              color: AppColors.darkGrey,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 16),
                  _AreaSection(
                    title: 'Récapitulatif',
                    subtitle: 'Vue rapide des informations saisies.',
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        _SummaryBlock(
                          title: 'Déclencheur',
                          serviceName: _actionService?.displayName ?? 'Non défini',
                          actionName: _selectedAction?.label ?? 'Aucun',
                          fields: _selectedAction?.fields ?? const [],
                          values: _actionFieldValues,
                          formatFieldValue: _formatFieldValue,
                        ),
                        const SizedBox(height: 12),
                        _SummaryBlock(
                          title: 'Réaction',
                          serviceName:
                              _reactionService?.displayName ?? 'Non défini',
                          actionName: _selectedReaction?.label ?? 'Aucune',
                          fields: _selectedReaction?.fields ?? const [],
                          values: _reactionFieldValues,
                          formatFieldValue: _formatFieldValue,
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            if (_createError != null)
              Padding(
                padding: const EdgeInsets.only(top: 16),
                child: Card(
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
              ),
            const SizedBox(height: 16),
            Row(
              children: [
                if (_wizardStep != AreaWizardStep.action)
                  Expanded(
                    child: OutlinedButton(
                      onPressed: _goToPreviousStep,
                      child: const Text('Étape précédente'),
                    ),
                  ),
                if (_wizardStep != AreaWizardStep.action)
                  const SizedBox(width: 12),
                Expanded(
                  child: ElevatedButton(
                    onPressed: _wizardStep == AreaWizardStep.details
                        ? (_isCreating ? null : _handleCreateArea)
                        : (_wizardStep == AreaWizardStep.action
                            ? (_canProceedAction ? _goToNextStep : null)
                            : (_canProceedReaction ? _goToNextStep : null)),
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
                        : const Text('Créer l\'AREA'))
                        : const Text('Continuer'),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class _WizardStepInfo {
  final AreaWizardStep step;
  final String title;
  final String description;

  const _WizardStepInfo({
    required this.step,
    required this.title,
    required this.description,
  });
}

class _AreaSection extends StatelessWidget {
  final String title;
  final String subtitle;
  final Widget child;

  const _AreaSection({
    required this.title,
    required this.subtitle,
    required this.child,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              title,
              style: theme.textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(height: 4),
            Text(
              subtitle,
              style: theme.textTheme.bodySmall?.copyWith(
                color: AppColors.darkGrey,
              ),
            ),
            const SizedBox(height: 12),
            child,
          ],
        ),
      ),
    );
  }
}

class _SectionLabel extends StatelessWidget {
  final String label;

  const _SectionLabel({required this.label});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Text(
        label,
        style: theme.textTheme.bodySmall?.copyWith(
          color: AppColors.darkGrey,
          fontWeight: FontWeight.w600,
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
  final String Function(DateTime) formatDate;

  const _FieldList({
    required this.keyPrefix,
    required this.fields,
    required this.values,
    required this.onChanged,
    required this.onDatePick,
    required this.formatDate,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Column(
      children: fields.map((field) {
        final value = values[field.name] ?? '';
        final label =
            field.required ? '${field.label} *' : field.label;
        if (field.type == AreaFieldType.date) {
          final display = value.isEmpty
              ? 'Sélectionner une date'
              : _formatSafeDate(formatDate, value);
          return Padding(
            padding: const EdgeInsets.only(bottom: 12),
            child: Material(
              color: Colors.transparent,
              child: InkWell(
                borderRadius: BorderRadius.circular(8),
                onTap: () => onDatePick(field.name, onChanged),
                child: InputDecorator(
                  decoration: InputDecoration(labelText: label),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        display,
                        style: theme.textTheme.bodyMedium?.copyWith(
                          color: value.isEmpty
                              ? AppColors.darkGrey
                              : AppColors.almostBlack,
                        ),
                      ),
                      const Icon(Icons.calendar_today, size: 18),
                    ],
                  ),
                ),
              ),
            ),
          );
        }

        final isDescription =
            field.name.toLowerCase().contains('description');
        return Padding(
          padding: const EdgeInsets.only(bottom: 12),
          child: TextFormField(
            key: ValueKey('$keyPrefix-${field.name}'),
            initialValue: value,
            keyboardType: field.type == AreaFieldType.number
                ? TextInputType.number
                : TextInputType.text,
            inputFormatters: field.type == AreaFieldType.number
                ? [FilteringTextInputFormatter.digitsOnly]
                : null,
            maxLines: isDescription ? 4 : 1,
            decoration: InputDecoration(labelText: label),
            onChanged: (newValue) => onChanged(field.name, newValue),
          ),
        );
      }).toList(),
    );
  }

  String _formatSafeDate(
    String Function(DateTime) formatDate,
    String value,
  ) {
    try {
      return formatDate(DateTime.parse(value));
    } catch (_) {
      return value;
    }
  }
}

class _SummaryBlock extends StatelessWidget {
  final String title;
  final String serviceName;
  final String actionName;
  final List<AreaFieldDefinition> fields;
  final Map<String, String> values;
  final String Function(AreaFieldDefinition, Map<String, String>)
      formatFieldValue;

  const _SummaryBlock({
    required this.title,
    required this.serviceName,
    required this.actionName,
    required this.fields,
    required this.values,
    required this.formatFieldValue,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.lightGrey,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.grey),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            title,
            style: theme.textTheme.bodySmall?.copyWith(
              color: AppColors.darkGrey,
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 6),
          Text(
            serviceName,
            style: theme.textTheme.titleSmall?.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          Text(
            actionName,
            style: theme.textTheme.bodySmall?.copyWith(
              color: AppColors.darkGrey,
            ),
          ),
          const SizedBox(height: 8),
          ...fields.map(
            (field) => Padding(
              padding: const EdgeInsets.only(bottom: 4),
              child: Text(
                '${field.label}: ${formatFieldValue(field, values)}',
                style: theme.textTheme.bodySmall?.copyWith(
                  color: AppColors.darkGrey,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
