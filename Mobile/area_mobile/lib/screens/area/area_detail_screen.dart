import 'package:flutter/material.dart';
import '../../models/area_definitions.dart';
import '../../theme/colors.dart';

class AreaDetailScreen extends StatelessWidget {
  final CreatedArea area;

  const AreaDetailScreen({
    super.key,
    required this.area,
  });

  String _formatDateTime(DateTime dt) {
    String two(int value) => value.toString().padLeft(2, '0');
    return '${dt.year}-${two(dt.month)}-${two(dt.day)} ${two(dt.hour)}:${two(dt.minute)}';
  }

  String _formatIso(String value) {
    try {
      return _formatDateTime(DateTime.parse(value));
    } catch (_) {
      return value;
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colors = context.appColors;
    final title = area.summary.isNotEmpty ? area.summary : area.name;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Détail de l\'AREA'),
      ),
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.all(20),
          children: [
            Container(
              padding: const EdgeInsets.all(18),
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(16),
                gradient: LinearGradient(
                  colors: [area.gradient.from, area.gradient.to],
                ),
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
                      _Badge(label: _badgeFrom(area.actionService)),
                      const SizedBox(width: 10),
                      const Icon(Icons.arrow_forward,
                          color: Colors.white, size: 18),
                      const SizedBox(width: 10),
                      _Badge(label: _badgeFrom(area.reactionService)),
                      const Spacer(),
                      _statusDot(),
                    ],
                  ),
                  const SizedBox(height: 14),
                  Text(
                    title,
                    style: theme.textTheme.titleLarge?.copyWith(
                      color: Colors.white,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  const SizedBox(height: 6),
                  Text(
                    area.name,
                    style: theme.textTheme.bodySmall?.copyWith(
                      color: Colors.white70,
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 16),
            _SectionCard(
              title: 'Service d\'action',
              children: [
                Text(
                  area.actionService,
                  style: theme.textTheme.titleSmall?.copyWith(
                    fontWeight: FontWeight.w600,
                  ),
                ),
                Text(
                  area.actionName,
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: colors.darkGrey,
                  ),
                ),
              ],
            ),
            _SectionCard(
              title: 'Service de réaction',
              children: [
                Text(
                  area.reactionService,
                  style: theme.textTheme.titleSmall?.copyWith(
                    fontWeight: FontWeight.w600,
                  ),
                ),
                Text(
                  area.reactionName,
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: colors.darkGrey,
                  ),
                ),
              ],
            ),
            _SectionCard(
              title: 'Paramètres',
              children: [
                _ParameterRow(
                  label: 'Début',
                  value: _formatIso(area.startTime),
                ),
                const SizedBox(height: 8),
                _ParameterRow(
                  label: 'Fin',
                  value: _formatIso(area.endTime),
                ),
                const SizedBox(height: 8),
                _ParameterRow(
                  label: 'Délai',
                  value: '${area.delay}s',
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

  Widget _statusDot() {
    return Container(
      width: 10,
      height: 10,
      decoration: BoxDecoration(
        color: Colors.greenAccent,
        shape: BoxShape.circle,
        boxShadow: [
          BoxShadow(
            color: Colors.greenAccent.withOpacity(0.4),
            blurRadius: 8,
            offset: const Offset(0, 2),
          ),
        ],
      ),
    );
  }
}

class _ParameterRow extends StatelessWidget {
  final String label;
  final String value;

  const _ParameterRow({
    required this.label,
    required this.value,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colors = context.appColors;
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          label,
          style: theme.textTheme.bodySmall?.copyWith(
            color: colors.darkGrey,
          ),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: Text(
            value,
            textAlign: TextAlign.right,
            style: theme.textTheme.bodySmall?.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
        ),
      ],
    );
  }
}

class _Badge extends StatelessWidget {
  final String label;

  const _Badge({required this.label});

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

class _SectionCard extends StatelessWidget {
  final String title;
  final List<Widget> children;

  const _SectionCard({required this.title, required this.children});

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
