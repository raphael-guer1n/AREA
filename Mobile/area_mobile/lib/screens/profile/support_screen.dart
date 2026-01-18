import 'package:flutter/material.dart';
import '../../theme/colors.dart';

class SupportScreen extends StatelessWidget {
  const SupportScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final colors = context.appColors;
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final surface = isDark ? colors.white.withOpacity(0.08) : Colors.white;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Support'),
        backgroundColor: Colors.transparent,
        foregroundColor: isDark ? Colors.white : colors.almostBlack,
        elevation: 0,
      ),
      backgroundColor: isDark ? colors.white : colors.lightGrey,
      body: Padding(
        padding: const EdgeInsets.all(20),
        child: Card(
          color: surface,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: const [
              ListTile(
                leading: Icon(Icons.forum_outlined),
                title: Text('Help Center'),
                subtitle: Text('FAQ, guides et support rapide'),
              ),
              Divider(height: 1),
              ListTile(
                leading: Icon(Icons.mail_outline),
                title: Text('Contact'),
                subtitle: Text('support@area.app'),
              ),
              Divider(height: 1),
              ListTile(
                leading: Icon(Icons.privacy_tip_outlined),
                title: Text('Confidentialité'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
