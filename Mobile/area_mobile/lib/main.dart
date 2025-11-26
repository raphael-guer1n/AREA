import 'package:flutter/material.dart';
import 'theme/app_theme.dart';
import 'screens/main_shell.dart';

void main() {
  runApp(const ActionReactionApp());
}

class ActionReactionApp extends StatelessWidget {
  const ActionReactionApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Action‑Reaction',
      debugShowCheckedModeBanner: false,
      theme: appTheme,
      home: const MainShell(),
    );
  }
}