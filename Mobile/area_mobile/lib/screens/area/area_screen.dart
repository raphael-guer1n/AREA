import 'package:flutter/material.dart';
import '../../models/area_definitions.dart';
import '../../theme/colors.dart';
import 'area_create_screen.dart';
import 'area_detail_screen.dart';

class AreaScreen extends StatefulWidget {
  const AreaScreen({super.key});

  @override
  State<AreaScreen> createState() => _AreaScreenState();
}

class _AreaScreenState extends State<AreaScreen> {
  final TextEditingController _searchController = TextEditingController();
  final List<CreatedArea> _areas = [];
  final Map<String, bool> _areaStatus = {};
  String _searchTerm = '';

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  Future<void> _openCreateArea() async {
    final created = await Navigator.of(context).push<CreatedArea>(
      MaterialPageRoute(builder: (_) => const CreateAreaScreen()),
    );
    if (created != null) {
      setState(() {
        _areas.insert(0, created);
        _areaStatus[created.id] = true;
      });
    }
  }

  void _openAreaDetail(CreatedArea area) {
    Navigator.of(context).push(
      MaterialPageRoute(builder: (_) => AreaDetailScreen(area: area)),
    );
  }

  void _openAreaPreview(CreatedArea area) {
    final colors = context.appColors;
    final theme = Theme.of(context);
    final isActive = _areaStatus[area.id] ?? true;
    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.transparent,
      builder: (context) {
        return Container(
          decoration: BoxDecoration(
            color: theme.scaffoldBackgroundColor,
            borderRadius: const BorderRadius.vertical(
              top: Radius.circular(18),
            ),
            boxShadow: [
              BoxShadow(
                color: colors.grey.withOpacity(0.2),
                blurRadius: 18,
                offset: const Offset(0, -6),
              ),
            ],
          ),
          padding: const EdgeInsets.fromLTRB(20, 14, 20, 24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Container(
                    width: 42,
                    height: 42,
                    decoration: BoxDecoration(
                      borderRadius: BorderRadius.circular(12),
                      gradient: LinearGradient(
                        colors: [area.gradient.from, area.gradient.to],
                      ),
                    ),
                    child: Center(
                      child: Text(
                        _badgeFrom(area.actionService),
                        style: const TextStyle(
                          color: Colors.white,
                          fontWeight: FontWeight.w700,
                        ),
                      ),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          area.summary.isNotEmpty
                              ? area.summary
                              : area.name,
                          style: theme.textTheme.titleMedium,
                        ),
                        Text(
                          '${area.actionName} → ${area.reactionName}',
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: colors.darkGrey,
                          ),
                        ),
                      ],
                    ),
                  ),
                  PopupMenuButton<String>(
                    onSelected: (value) {
                      Navigator.of(context).pop();
                      if (value == 'toggle') {
                        _toggleArea(area);
                      } else if (value == 'delete') {
                        _deleteArea(area);
                      } else if (value == 'detail') {
                        _openAreaDetail(area);
                      }
                    },
                    itemBuilder: (context) => [
                      PopupMenuItem(
                        value: 'toggle',
                        child: Text(isActive ? 'Désactiver' : 'Activer'),
                      ),
                      const PopupMenuItem(
                        value: 'detail',
                        child: Text('Voir le détail'),
                      ),
                      const PopupMenuItem(
                        value: 'delete',
                        child: Text('Supprimer'),
                      ),
                    ],
                    icon: const Icon(Icons.more_horiz),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              Text(
                area.summary.isNotEmpty ? area.summary : area.name,
                style: theme.textTheme.bodyMedium,
              ),
              const SizedBox(height: 8),
              Text(
                'Déclencheur: ${area.actionName}\nRéaction: ${area.reactionName}',
                style: theme.textTheme.bodySmall?.copyWith(
                  color: colors.darkGrey,
                ),
              ),
            ],
          ),
        );
      },
    );
  }

  void _toggleArea(CreatedArea area) {
    setState(() {
      final current = _areaStatus[area.id] ?? true;
      _areaStatus[area.id] = !current;
    });
  }

  void _deleteArea(CreatedArea area) {
    setState(() {
      _areas.removeWhere((a) => a.id == area.id);
      _areaStatus.remove(area.id);
    });
  }

  List<CreatedArea> get _filteredAreas {
    final term = _searchTerm.trim().toLowerCase();
    if (term.isEmpty) return _areas;
    return _areas.where((area) {
      final haystack = [
        area.name,
        area.summary,
        area.actionService,
        area.reactionService,
        area.actionName,
        area.reactionName,
      ].join(' ').toLowerCase();
      return haystack.contains(term);
    }).toList();
  }

  String _badgeFrom(String value) {
    if (value.isEmpty) return '--';
    final trimmed = value.trim();
    if (trimmed.length <= 2) return trimmed.toUpperCase();
    return trimmed.substring(0, 2).toUpperCase();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colors = context.appColors;
    const horizontalPadding = 20.0;
    final activeCount = _areas.length;
    final totalCount = _areas.length;
    final filteredAreas = _filteredAreas;

    return Scaffold(
      backgroundColor: theme.colorScheme.background,
      body: SafeArea(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(
                horizontalPadding,
                18,
                horizontalPadding,
                12,
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'AREA',
                    style: theme.textTheme.displayLarge,
                  ),
                  const SizedBox(height: 8),
                  Text(
                    'Créez vos automatisations en quelques étapes.',
                    style: theme.textTheme.bodyMedium?.copyWith(
                      color: colors.darkGrey,
                    ),
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      Expanded(
                        child: Card(
                          child: Padding(
                            padding: const EdgeInsets.all(14),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  'Areas actives',
                                  style: theme.textTheme.bodySmall?.copyWith(
                                    color: colors.darkGrey,
                                  ),
                                ),
                                const SizedBox(height: 6),
                                Text(
                                  activeCount.toString(),
                                  style: theme.textTheme.headlineMedium?.copyWith(
                                    color: colors.midBlue,
                                    fontWeight: FontWeight.w700,
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Card(
                          child: Padding(
                            padding: const EdgeInsets.all(14),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  'Areas créées',
                                  style: theme.textTheme.bodySmall?.copyWith(
                                    color: colors.darkGrey,
                                  ),
                                ),
                                const SizedBox(height: 6),
                                Text(
                                  totalCount.toString(),
                                  style: theme.textTheme.headlineMedium?.copyWith(
                                    color: colors.midBlue,
                                    fontWeight: FontWeight.w700,
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  SizedBox(
                    width: double.infinity,
                    height: 48,
                    child: ElevatedButton.icon(
                      onPressed: _openCreateArea,
                      icon: const Icon(Icons.add),
                      label: const Text('Créer une AREA'),
                    ),
                  ),
                ],
              ),
            ),
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: horizontalPadding),
              child: TextField(
                controller: _searchController,
                onChanged: (value) {
                  setState(() {
                    _searchTerm = value;
                  });
                },
                decoration: InputDecoration(
                  hintText: 'Rechercher une area...',
                  prefixIcon: const Icon(Icons.search),
                  suffixIcon: _searchTerm.isNotEmpty
                      ? IconButton(
                          icon: const Icon(Icons.clear),
                          onPressed: () {
                            _searchController.clear();
                            setState(() {
                              _searchTerm = '';
                            });
                          },
                        )
                      : null,
                ),
              ),
            ),
            const SizedBox(height: 16),
            Expanded(
              child: _areas.isEmpty
                  ? _EmptyAreaState(onCreate: _openCreateArea)
                  : filteredAreas.isEmpty
                      ? _NoResultsState(
                          onReset: () {
                            _searchController.clear();
                            setState(() {
                              _searchTerm = '';
                            });
                          },
                        )
                      : GridView.builder(
                          padding: const EdgeInsets.fromLTRB(
                            horizontalPadding,
                            0,
                            horizontalPadding,
                            20,
                          ),
                          gridDelegate:
                              const SliverGridDelegateWithFixedCrossAxisCount(
                            crossAxisCount: 2,
                            childAspectRatio: 1.1,
                            crossAxisSpacing: 12,
                            mainAxisSpacing: 12,
                          ),
                          itemCount: filteredAreas.length,
                                  itemBuilder: (context, index) {
                                    final area = filteredAreas[index];
                                    return _AreaCard(
                                      area: area,
                                      isActive:
                                          _areaStatus[area.id] ?? true,
                                      onPreview: () => _openAreaPreview(area),
                                      onOpenDetail: () => _openAreaDetail(area),
                                      onToggle: () => _toggleArea(area),
                                      onDelete: () => _deleteArea(area),
                                    );
                                  },
                                ),
            ),
          ],
        ),
      ),
    );
  }
}

class _AreaCard extends StatelessWidget {
  final CreatedArea area;
  final bool isActive;
  final VoidCallback onPreview;
  final VoidCallback onOpenDetail;
  final VoidCallback onToggle;
  final VoidCallback onDelete;

  const _AreaCard({
    required this.area,
    required this.isActive,
    required this.onPreview,
    required this.onOpenDetail,
    required this.onToggle,
    required this.onDelete,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colors = context.appColors;
    final actionBadge = _badgeFrom(area.actionService);
    final reactionBadge = _badgeFrom(area.reactionService);
    final title = area.summary.isNotEmpty ? area.summary : area.name;

    return InkWell(
      onTap: onPreview,
      borderRadius: BorderRadius.circular(16),
      child: Container(
        padding: const EdgeInsets.all(14),
        decoration: BoxDecoration(
          gradient: LinearGradient(
            colors: [area.gradient.from, area.gradient.to],
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
          ),
          borderRadius: BorderRadius.circular(16),
          boxShadow: [
            BoxShadow(
              color: colors.grey.withOpacity(0.2),
              blurRadius: 14,
              offset: const Offset(0, 6),
            ),
          ],
        ),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                _ServiceBadge(label: actionBadge),
                const SizedBox(width: 8),
                const Icon(Icons.arrow_forward, size: 18, color: Colors.white),
                const SizedBox(width: 8),
                _ServiceBadge(label: reactionBadge),
                const Spacer(),
                _statusDot(isActive),
                PopupMenuButton<String>(
                  onSelected: (value) {
                    if (value == 'toggle') {
                      onToggle();
                    } else if (value == 'delete') {
                      onDelete();
                    } else if (value == 'detail') {
                      onOpenDetail();
                    }
                  },
                  icon: const Icon(Icons.more_horiz, color: Colors.white),
                  itemBuilder: (context) => [
                    PopupMenuItem(
                      value: 'toggle',
                      child: Text(isActive ? 'Désactiver' : 'Activer'),
                    ),
                    const PopupMenuItem(
                      value: 'detail',
                      child: Text('Voir le détail'),
                    ),
                    const PopupMenuItem(
                      value: 'delete',
                      child: Text('Supprimer'),
                    ),
                  ],
                ),
              ],
            ),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  title,
                  style: theme.textTheme.titleMedium?.copyWith(
                    color: Colors.white,
                    fontWeight: FontWeight.w700,
                  ),
                ),
                const SizedBox(height: 6),
                Text(
                  '${area.actionName} → ${area.reactionName}',
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: Colors.white70,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  String _badgeFrom(String value) {
    if (value.isEmpty) return '--';
    final trimmed = value.trim();
    if (trimmed.length <= 2) return trimmed.toUpperCase();
    return trimmed.substring(0, 2).toUpperCase();
  }

  Widget _statusDot(bool isActive) {
    final color = isActive ? Colors.greenAccent : Colors.redAccent;
    return Container(
      width: 10,
      height: 10,
      decoration: BoxDecoration(
        color: color,
        shape: BoxShape.circle,
        boxShadow: [
          BoxShadow(
            color: color.withOpacity(0.4),
            blurRadius: 8,
            offset: const Offset(0, 2),
          ),
        ],
      ),
    );
  }
}

class _ServiceBadge extends StatelessWidget {
  final String label;

  const _ServiceBadge({required this.label});

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 34,
      height: 34,
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.18),
        borderRadius: BorderRadius.circular(10),
        border: Border.all(color: Colors.white.withOpacity(0.25)),
      ),
      child: Center(
        child: Text(
          label,
          style: const TextStyle(
            fontWeight: FontWeight.w700,
            color: Colors.white,
          ),
        ),
      ),
    );
  }
}

class _EmptyAreaState extends StatelessWidget {
  final VoidCallback onCreate;

  const _EmptyAreaState({required this.onCreate});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colors = context.appColors;
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              Icons.auto_awesome,
              size: 56,
              color: colors.darkGrey,
            ),
            const SizedBox(height: 16),
            Text(
              'Pas encore d\'AREA',
              style: theme.textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            Text(
              'Créez votre première automation pour la voir apparaître ici.',
              textAlign: TextAlign.center,
              style: theme.textTheme.bodySmall?.copyWith(
                color: colors.darkGrey,
              ),
            ),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: onCreate,
              child: const Text('Créer une AREA'),
            ),
          ],
        ),
      ),
    );
  }
}

class _NoResultsState extends StatelessWidget {
  final VoidCallback onReset;

  const _NoResultsState({required this.onReset});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colors = context.appColors;
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              Icons.search_off,
              size: 56,
              color: colors.darkGrey,
            ),
            const SizedBox(height: 16),
            Text(
              'Aucune area trouvée',
              style: theme.textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            Text(
              'Essayez un autre mot-clé ou réinitialisez la recherche.',
              textAlign: TextAlign.center,
              style: theme.textTheme.bodySmall?.copyWith(
                color: colors.darkGrey,
              ),
            ),
            const SizedBox(height: 16),
            OutlinedButton(
              onPressed: onReset,
              child: const Text('Réinitialiser'),
            ),
          ],
        ),
      ),
    );
  }
}
