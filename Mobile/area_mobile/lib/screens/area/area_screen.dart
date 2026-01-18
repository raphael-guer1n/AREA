import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../models/area_backend_models.dart';
import '../../providers/area_provider.dart';
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
  String _searchTerm = '';

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<AreaProvider>().loadAreas();
    });
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  Future<void> _openCreateArea() async {
    final created = await Navigator.of(context).push<bool>(
      MaterialPageRoute(builder: (_) => const CreateAreaScreen()),
    );

    if (created == true && mounted) {
      await context.read<AreaProvider>().loadAreas();
    }
  }

  void _openAreaDetail(AreaDto area) {
    Navigator.of(context).push(
      MaterialPageRoute(builder: (_) => AreaDetailScreen(area: area)),
    );
  }

  List<AreaDto> _filteredAreas(List<AreaDto> areas) {
    final term = _searchTerm.trim().toLowerCase();
    if (term.isEmpty) return areas;

    return areas.where((area) {
      final action = area.actions.isNotEmpty ? area.actions.first : null;
      final reaction = area.reactions.isNotEmpty ? area.reactions.first : null;

      final haystack = [
        area.name,
        action?.service ?? '',
        action?.title ?? '',
        reaction?.service ?? '',
        reaction?.title ?? '',
      ].join(' ').toLowerCase();

      return haystack.contains(term);
    }).toList();
  }

  Future<void> _toggleArea(AreaDto area) async {
    final provider = context.read<AreaProvider>();
    await provider.toggleArea(area);

    if (!mounted) return;

    final err = provider.error;
    if (err != null) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(err), backgroundColor: Colors.red),
      );
    }
  }

  Future<void> _confirmAndDeleteArea(AreaDto area) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Supprimer l’AREA'),
        content: Text(
          'Voulez-vous vraiment supprimer "${area.name}" ?\n'
          'Cette action est irréversible.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('Annuler'),
          ),
          ElevatedButton(
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.red.shade600,
              foregroundColor: Colors.white,
            ),
            onPressed: () => Navigator.of(context).pop(true),
            child: const Text('Supprimer'),
          ),
        ],
      ),
    );

    if (confirmed != true) return;

    final provider = context.read<AreaProvider>();
    await provider.deleteArea(area.id);

    if (!mounted) return;

    final err = provider.error;
    if (err != null) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(err), backgroundColor: Colors.red),
      );
      return;
    }

    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('AREA supprimée'),
        backgroundColor: Colors.green,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colors = context.appColors;

    final provider = context.watch<AreaProvider>();
    final areas = provider.areas;
    final filtered = _filteredAreas(areas);

    final activeCount = areas.where((a) => a.active).length;
    final totalCount = areas.length;

    return Scaffold(
      backgroundColor: theme.colorScheme.background,
      body: SafeArea(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(20, 18, 20, 12),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('AREA', style: theme.textTheme.displayLarge),
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
                            padding: const EdgeInsets.symmetric(
                              horizontal: 10,
                              vertical: 8,
                            ),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  'Areas actives',
                                  style: theme.textTheme.bodySmall?.copyWith(
                                    color: colors.darkGrey,
                                  ),
                                ),
                                const SizedBox(height: 2),
                                Text(
                                  activeCount.toString(),
                                  style: theme.textTheme.titleMedium?.copyWith(
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
                            padding: const EdgeInsets.symmetric(
                              horizontal: 10,
                              vertical: 8,
                            ),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  'Areas créées',
                                  style: theme.textTheme.bodySmall?.copyWith(
                                    color: colors.darkGrey,
                                  ),
                                ),
                                const SizedBox(height: 2),
                                Text(
                                  totalCount.toString(),
                                  style: theme.textTheme.titleMedium?.copyWith(
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
              padding: const EdgeInsets.symmetric(horizontal: 20),
              child: TextField(
                controller: _searchController,
                onChanged: (value) => setState(() => _searchTerm = value),
                decoration: InputDecoration(
                  hintText: 'Rechercher une area...',
                  prefixIcon: const Icon(Icons.search),
                  suffixIcon: _searchTerm.isNotEmpty
                      ? IconButton(
                          icon: const Icon(Icons.clear),
                          onPressed: () {
                            _searchController.clear();
                            setState(() => _searchTerm = '');
                          },
                        )
                      : null,
                ),
              ),
            ),
            const SizedBox(height: 16),
            Expanded(
              child: provider.isLoading
                  ? const Center(child: CircularProgressIndicator())
                  : provider.error != null
                      ? _ErrorState(
                          error: provider.error!,
                          onRetry: () => context.read<AreaProvider>().loadAreas(),
                        )
                      : RefreshIndicator(
                          onRefresh: () => context.read<AreaProvider>().loadAreas(),
                          child: filtered.isEmpty
                              ? _EmptyState(onCreate: _openCreateArea)
                              : GridView.builder(
                                  padding: const EdgeInsets.fromLTRB(
                                    20,
                                    0,
                                    20,
                                    20,
                                  ),
                                  gridDelegate:
                                      const SliverGridDelegateWithFixedCrossAxisCount(
                                    crossAxisCount: 2,
                                    childAspectRatio: 1.1,
                                    crossAxisSpacing: 12,
                                    mainAxisSpacing: 12,
                                  ),
                                  itemCount: filtered.length,
                                  itemBuilder: (context, index) {
                                    final area = filtered[index];
                                    return _AreaCard(
                                      area: area,
                                      onOpenDetail: () => _openAreaDetail(area),
                                      onToggle: () => _toggleArea(area),
                                      onDelete: () => _confirmAndDeleteArea(area),
                                    );
                                  },
                                ),
                        ),
            ),
          ],
        ),
      ),
    );
  }
}

class _AreaCard extends StatelessWidget {
  final AreaDto area;
  final VoidCallback onOpenDetail;
  final VoidCallback onToggle;
  final VoidCallback onDelete;

  const _AreaCard({
    required this.area,
    required this.onOpenDetail,
    required this.onToggle,
    required this.onDelete,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colors = context.appColors;

    final title = area.name;
    final gradient = _gradientFor(area.id);
    final badge = _buildBadge(title);

    return InkWell(
      onTap: onOpenDetail,
      borderRadius: BorderRadius.circular(16),
      child: Container(
        padding: const EdgeInsets.all(14),
        decoration: BoxDecoration(
          gradient: LinearGradient(
            colors: gradient,
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
          ),
          borderRadius: BorderRadius.circular(16),
          boxShadow: [
            BoxShadow(
              color: colors.grey.withOpacity(0.25),
              blurRadius: 16,
              offset: const Offset(0, 8),
            ),
          ],
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                CircleAvatar(
                  radius: 18,
                  backgroundColor: Colors.white.withOpacity(0.18),
                  child: Text(
                    badge,
                    style: const TextStyle(
                      color: Colors.white,
                      fontWeight: FontWeight.w700,
                      fontSize: 14,
                    ),
                  ),
                ),
                const Spacer(),
                Container(
                  width: 30,
                  height: 30,
                  decoration: BoxDecoration(
                    color: Colors.white.withOpacity(0.16),
                    shape: BoxShape.circle,
                    border: Border.all(color: Colors.white.withOpacity(0.35)),
                  ),
                  child: PopupMenuButton<String>(
                    padding: EdgeInsets.zero,
                    onSelected: (value) {
                      if (value == 'toggle') {
                        onToggle();
                      } else if (value == 'delete') {
                        onDelete();
                      } else if (value == 'detail') {
                        onOpenDetail();
                      }
                    },
                    icon: const Icon(Icons.more_horiz, color: Colors.white, size: 18),
                    itemBuilder: (context) => [
                      PopupMenuItem(
                        value: 'toggle',
                        child: Text(area.active ? 'Désactiver' : 'Activer'),
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
                ),
              ],
            ),
            const SizedBox(height: 12),
            Text(
              title,
              maxLines: 2,
              overflow: TextOverflow.ellipsis,
              style: theme.textTheme.titleMedium?.copyWith(
                color: Colors.white,
                fontSize: 15,
                fontWeight: FontWeight.w600,
              ),
            ),
            const Spacer(),
            Align(
              alignment: Alignment.bottomLeft,
              child: _statusDot(area.active),
            ),
          ],
        ),
      ),
    );
  }

  String _buildBadge(String name) {
    final parts = name.trim().split(' ').where((part) => part.isNotEmpty).toList();
    if (parts.length >= 2) {
      return '${parts[0][0]}${parts[1][0]}'.toUpperCase();
    }
    return name.isNotEmpty
        ? name.substring(0, name.length >= 2 ? 2 : 1).toUpperCase()
        : 'A';
  }

  List<Color> _gradientFor(int id) {
    const palette = [
      [Color(0xFF002642), Color(0xFF0B3C5D)],
      [Color(0xFF840032), Color(0xFFA33A60)],
      [Color(0xFFE59500), Color(0xFFF2B344)],
      [Color(0xFF5B834D), Color(0xFF68915A)],
      [Color(0xFF02040F), Color(0xFF1B2640)],
    ];
    return palette[id.abs() % palette.length];
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

class _EmptyState extends StatelessWidget {
  final VoidCallback onCreate;

  const _EmptyState({required this.onCreate});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colors = context.appColors;

    return ListView(
      physics: const AlwaysScrollableScrollPhysics(),
      children: [
        const SizedBox(height: 60),
        Center(
          child: Padding(
            padding: const EdgeInsets.all(20),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(Icons.auto_awesome, size: 56, color: colors.darkGrey),
                const SizedBox(height: 16),
                Text("Pas encore d'AREA", style: theme.textTheme.titleMedium),
                const SizedBox(height: 8),
                Text(
                  "Créez votre première automation pour la voir apparaître ici.",
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
        ),
      ],
    );
  }
}

class _ErrorState extends StatelessWidget {
  final String error;
  final VoidCallback onRetry;

  const _ErrorState({
    required this.error,
    required this.onRetry,
  });

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
            Icon(Icons.error_outline, size: 56, color: colors.darkGrey),
            const SizedBox(height: 12),
            Text('Erreur', style: theme.textTheme.titleMedium),
            const SizedBox(height: 8),
            Text(error, textAlign: TextAlign.center),
            const SizedBox(height: 12),
            ElevatedButton.icon(
              onPressed: onRetry,
              icon: const Icon(Icons.refresh),
              label: const Text('Réessayer'),
            ),
          ],
        ),
      ),
    );
  }
}
