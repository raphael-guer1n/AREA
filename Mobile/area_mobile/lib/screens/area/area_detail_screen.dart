import 'package:flutter/material.dart';
import '../../models/area_backend_models.dart';
import '../../theme/colors.dart';

class AreaDetailScreen extends StatelessWidget {
  final AreaDto area;

  const AreaDetailScreen({
    super.key,
    required this.area,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colors = context.appColors;

    final action = area.actions.isNotEmpty ? area.actions.first : null;
    final reaction = area.reactions.isNotEmpty ? area.reactions.first : null;

    return Scaffold(
      appBar: AppBar(
        title: const Text("Détail de l'AREA"),
      ),
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.all(20),
          children: [
            Text(
              area.name,
              style: theme.textTheme.titleLarge,
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                _statusDot(area.active),
                const SizedBox(width: 8),
                Text(
                  area.active ? 'Active' : 'Inactive',
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: area.active ? Colors.green.shade700 : colors.darkGrey,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            _SectionCard(
              title: 'Action',
              children: [
                Text(action?.service ?? '—'),
                Text(
                  action?.title ?? '—',
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: colors.darkGrey,
                  ),
                ),
              ],
            ),
            _SectionCard(
              title: 'Réaction',
              children: [
                Text(reaction?.service ?? '—'),
                Text(
                  reaction?.title ?? '—',
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: colors.darkGrey,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _statusDot(bool isActive) {
    final color = isActive ? Colors.greenAccent : Colors.redAccent;
    return Container(
      width: 10,
      height: 10,
      decoration: BoxDecoration(
        color: color,
        shape: BoxShape.circle,
      ),
    );
  }
}

class _SectionCard extends StatelessWidget {
  final String title;
  final List<Widget> children;

  const _SectionCard({
    required this.title,
    required this.children,
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
          Text(
            title,
            style: theme.textTheme.titleSmall?.copyWith(
              fontWeight: FontWeight.w700,
            ),
          ),
          const SizedBox(height: 10),
          ...children,
        ],
      ),
    );
  }
}