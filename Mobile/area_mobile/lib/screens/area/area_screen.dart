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
      });
    }
  }

  void _openAreaDetail(CreatedArea area) {
    Navigator.of(context).push(
      MaterialPageRoute(builder: (_) => AreaDetailScreen(area: area)),
    );
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

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
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
              padding: const EdgeInsets.all(24.0),
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
                      color: AppColors.darkGrey,
                    ),
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      Expanded(
                        child: Card(
                          child: Padding(
                            padding: const EdgeInsets.all(16),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  'Areas actives',
                                  style: theme.textTheme.bodySmall?.copyWith(
                                    color: AppColors.darkGrey,
                                  ),
                                ),
                                const SizedBox(height: 6),
                                Text(
                                  activeCount.toString(),
                                  style: theme.textTheme.headlineMedium?.copyWith(
                                    color: AppColors.midBlue,
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
                            padding: const EdgeInsets.all(16),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  'Areas créées',
                                  style: theme.textTheme.bodySmall?.copyWith(
                                    color: AppColors.darkGrey,
                                  ),
                                ),
                                const SizedBox(height: 6),
                                Text(
                                  totalCount.toString(),
                                  style: theme.textTheme.headlineMedium?.copyWith(
                                    color: AppColors.midBlue,
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
              padding: const EdgeInsets.symmetric(horizontal: 24.0),
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
                      : ListView.builder(
                          padding: const EdgeInsets.only(
                            left: 24,
                            right: 24,
                            bottom: 24,
                          ),
                          itemCount: filteredAreas.length,
                          itemBuilder: (context, index) {
                            final area = filteredAreas[index];
                            return _AreaCard(
                              area: area,
                              onTap: () => _openAreaDetail(area),
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
  final VoidCallback onTap;

  const _AreaCard({
    required this.area,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final actionBadge = _badgeFrom(area.actionService);
    final reactionBadge = _badgeFrom(area.reactionService);
    final title = area.summary.isNotEmpty ? area.summary : area.name;

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Container(
              height: 6,
              decoration: BoxDecoration(
                borderRadius: const BorderRadius.vertical(
                  top: Radius.circular(12),
                ),
                gradient: LinearGradient(
                  colors: [area.gradient.from, area.gradient.to],
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      _ServiceBadge(label: actionBadge),
                      const SizedBox(width: 8),
                      const Icon(Icons.arrow_forward, size: 18),
                      const SizedBox(width: 8),
                      _ServiceBadge(label: reactionBadge),
                      const Spacer(),
                      Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 8,
                          vertical: 4,
                        ),
                        decoration: BoxDecoration(
                          color: Colors.green.shade50,
                          borderRadius: BorderRadius.circular(6),
                          border: Border.all(color: Colors.green.shade200),
                        ),
                        child: Text(
                          'Actif',
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: Colors.green.shade700,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  Text(
                    title,
                    style: theme.textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  const SizedBox(height: 6),
                  Text(
                    '${area.actionName} → ${area.reactionName}',
                    style: theme.textTheme.bodySmall?.copyWith(
                      color: AppColors.darkGrey,
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

  String _badgeFrom(String value) {
    if (value.isEmpty) return '--';
    final trimmed = value.trim();
    if (trimmed.length <= 2) return trimmed.toUpperCase();
    return trimmed.substring(0, 2).toUpperCase();
  }
}

class _ServiceBadge extends StatelessWidget {
  final String label;

  const _ServiceBadge({required this.label});

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 36,
      height: 36,
      decoration: BoxDecoration(
        color: AppColors.lightGrey,
        borderRadius: BorderRadius.circular(10),
      ),
      child: Center(
        child: Text(
          label,
          style: const TextStyle(
            fontWeight: FontWeight.w700,
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
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(
              Icons.auto_awesome,
              size: 56,
              color: AppColors.darkGrey,
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
                color: AppColors.darkGrey,
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
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(
              Icons.search_off,
              size: 56,
              color: AppColors.darkGrey,
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
                color: AppColors.darkGrey,
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
