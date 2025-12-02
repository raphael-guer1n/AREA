import 'package:area_mobile/screens/main_shell.dart';
import 'package:flutter/material.dart';
import 'theme/theme.dart';

void main() => runApp(const AreaApp());

class AreaApp extends StatelessWidget {
  const AreaApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'AREA',
      theme: areaTheme,
      debugShowCheckedModeBanner: false,
      home: const MainShell(),
    );
  }
}
