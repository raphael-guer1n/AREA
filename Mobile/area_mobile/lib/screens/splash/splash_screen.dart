import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../providers/auth_provider.dart';

class SplashScreen extends StatefulWidget {
  const SplashScreen({super.key});

  @override
  State<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends State<SplashScreen> {
  late final AuthProvider _authProvider;
  bool _delayDone = false;
  bool _navigated = false;

  @override
  void initState() {
    super.initState();
    _authProvider = context.read<AuthProvider>();
    _authProvider.addListener(_handleAuthChange);
    _startDelay();
  }

  @override
  void dispose() {
    _authProvider.removeListener(_handleAuthChange);
    super.dispose();
  }

  Future<void> _startDelay() async {
    await Future.delayed(const Duration(seconds: 3));
    if (!mounted) return;
    setState(() {
      _delayDone = true;
    });
    _attemptNavigation();
  }

  void _handleAuthChange() {
    _attemptNavigation();
  }

  void _attemptNavigation() {
    if (_navigated || !_delayDone || _authProvider.isLoading) return;
    _navigated = true;
    final target = _authProvider.isAuthenticated ? '/main' : '/login';
    Navigator.of(context).pushReplacementNamed(target);
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Scaffold(
      backgroundColor: theme.colorScheme.background,
      body: Center(
        child: Image.asset(
          'assets/logo.png',
          height: 110,
          fit: BoxFit.contain,
        ),
      ),
    );
  }
}
