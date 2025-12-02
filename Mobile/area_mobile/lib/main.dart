import 'package:flutter/material.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'screens/main_shell.dart';
import 'theme/theme.dart';

Future<void> main() async {

  await dotenv.load(fileName: ".env");
  runApp(const AreaApp());
}

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