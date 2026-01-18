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

    final actions = area.actions;
    final reactions = area.reactions;
    final actionSummary = actions.isEmpty
        ? '—'
        : '${actions.first.service} · ${actions.first.title}';
    final reactionSummary = reactions.isEmpty
        ? '—'
        : reactions
            .map((reaction) => '${reaction.service} · ${reaction.title}')
            .join('\n');

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
            const SizedBox(height: 16),
            _SectionCard(
              title: 'Résumé complet',
              children: [
                _summaryRow(
                  context,
                  label: 'Action',
                  value: actionSummary,
                ),
                _summaryRow(
                  context,
                  label: 'Réactions',
                  value: reactionSummary,
                ),
                const SizedBox(height: 8),
                Text(
                  'Actions',
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: colors.darkGrey,
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(height: 6),
                ...actions.isEmpty
                    ? [Text('—', style: theme.textTheme.bodySmall)]
                    : actions
                        .asMap()
                        .entries
                        .map(
                          (entry) => _ActionReactionCard(
                            index: entry.key + 1,
                            title: entry.value.title,
                            service: entry.value.service,
                            inputs: entry.value.input,
                          ),
                        )
                        .toList(),
                const SizedBox(height: 8),
                Text(
                  'Réactions',
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: colors.darkGrey,
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(height: 6),
                ...reactions.isEmpty
                    ? [Text('—', style: theme.textTheme.bodySmall)]
                    : reactions
                        .asMap()
                        .entries
                        .map(
                          (entry) => _ActionReactionCard(
                            index: entry.key + 1,
                            title: entry.value.title,
                            service: entry.value.service,
                            inputs: entry.value.input,
                          ),
                        )
                        .toList(),
              ],
            ),
            const SizedBox(height: 12),
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

  Widget _summaryRow(
    BuildContext context, {
    required String label,
    required String value,
  }) {
    final colors = context.appColors;

    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 80,
            child: Text(
              label,
              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: colors.darkGrey,
                    fontWeight: FontWeight.w600,
                  ),
            ),
          ),
          Expanded(
            child: Text(
              value,
              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: colors.almostBlack,
                  ),
            ),
          ),
        ],
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

class _ActionReactionCard extends StatelessWidget {
  final int index;
  final String title;
  final String service;
  final List<InputFieldDto> inputs;

  const _ActionReactionCard({
    required this.index,
    required this.title,
    required this.service,
    required this.inputs,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colors = context.appColors;

    return Container(
      margin: const EdgeInsets.only(bottom: 10),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: colors.lightGrey,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: colors.grey),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Étape $index',
            style: theme.textTheme.bodySmall?.copyWith(
              color: colors.darkGrey,
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            '$service · $title',
            style: theme.textTheme.bodyMedium?.copyWith(
              color: colors.almostBlack,
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 8),
          if (inputs.isEmpty)
            Text('Aucun paramètre', style: theme.textTheme.bodySmall)
          else
            Column(
              children: inputs.map((input) {
                return Padding(
                  padding: const EdgeInsets.only(bottom: 4),
                  child: Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      SizedBox(
                        width: 90,
                        child: Text(
                          input.name,
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: colors.darkGrey,
                          ),
                        ),
                      ),
                      Expanded(
                        child: Text(
                          input.value.isEmpty ? '—' : input.value,
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
      ),
    );
  }
}
