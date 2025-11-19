import 'package:flutter/material.dart';

void main() {
  runApp(const MyApp());
}

class AppColors {
  const AppColors._();

  static const background = Color(0xFFF4F4F5);
  static const foreground = Color(0xFF1B1B1F);
  static const muted = Color(0xFF4B5563);
  static const surface = Color(0xFFFFFFFF);
  static const surfaceBorder = Color(0xFFE4E4E7);
  static const accent = Color(0xFF111111);

  static const backgroundDark = Color(0xFF050505);
  static const foregroundDark = Color(0xFFF4F4F5);
  static const mutedDark = Color(0xFFA1A1AA);
  static const surfaceDark = Color(0xFF18181B);
  static const surfaceBorderDark = Color(0xFF27272A);
  static const accentDark = Color(0xFFF4F4F5);
}

Color mutedColor(BuildContext context) {
  return Theme.of(context).brightness == Brightness.dark
      ? AppColors.mutedDark
      : AppColors.muted;
}

Color surfaceBorderColor(BuildContext context) {
  return Theme.of(context).brightness == Brightness.dark
      ? AppColors.surfaceBorderDark
      : AppColors.surfaceBorder;
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'AREA',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        useMaterial3: true,
        scaffoldBackgroundColor: AppColors.background,
        colorScheme: const ColorScheme.light(
          background: AppColors.background,
          surface: AppColors.surface,
          primary: AppColors.accent,
          secondary: AppColors.foreground,
          onBackground: AppColors.foreground,
          onSurface: AppColors.foreground,
          onPrimary: AppColors.background,
          onSecondary: AppColors.background,
        ),
        textTheme: ThemeData.light().textTheme.apply(
              bodyColor: AppColors.foreground,
              displayColor: AppColors.foreground,
            ),
        appBarTheme: const AppBarTheme(
          backgroundColor: AppColors.background,
          foregroundColor: AppColors.foreground,
          elevation: 0,
        ),
        inputDecorationTheme: const InputDecorationTheme(
          filled: true,
          fillColor: AppColors.surface,
          border: OutlineInputBorder(
            borderSide: BorderSide(color: AppColors.surfaceBorder),
            borderRadius: BorderRadius.all(Radius.circular(16)),
          ),
          enabledBorder: OutlineInputBorder(
            borderSide: BorderSide(color: AppColors.surfaceBorder),
            borderRadius: BorderRadius.all(Radius.circular(16)),
          ),
          focusedBorder: OutlineInputBorder(
            borderSide: BorderSide(color: AppColors.foreground),
            borderRadius: BorderRadius.all(Radius.circular(16)),
          ),
          labelStyle: TextStyle(color: AppColors.foreground),
        ),
      ),
      darkTheme: ThemeData(
        useMaterial3: true,
        brightness: Brightness.dark,
        scaffoldBackgroundColor: AppColors.backgroundDark,
        colorScheme: const ColorScheme.dark(
          background: AppColors.backgroundDark,
          surface: AppColors.surfaceDark,
          primary: AppColors.accentDark,
          secondary: AppColors.foregroundDark,
          onBackground: AppColors.foregroundDark,
          onSurface: AppColors.foregroundDark,
          onPrimary: AppColors.backgroundDark,
          onSecondary: AppColors.backgroundDark,
        ),
        textTheme: ThemeData.dark().textTheme.apply(
              bodyColor: AppColors.foregroundDark,
              displayColor: AppColors.foregroundDark,
            ),
        inputDecorationTheme: const InputDecorationTheme(
          filled: true,
          fillColor: AppColors.surfaceDark,
          border: OutlineInputBorder(
            borderSide: BorderSide(color: AppColors.surfaceBorderDark),
            borderRadius: BorderRadius.all(Radius.circular(16)),
          ),
          enabledBorder: OutlineInputBorder(
            borderSide: BorderSide(color: AppColors.surfaceBorderDark),
            borderRadius: BorderRadius.all(Radius.circular(16)),
          ),
          focusedBorder: OutlineInputBorder(
            borderSide: BorderSide(color: AppColors.foregroundDark),
            borderRadius: BorderRadius.all(Radius.circular(16)),
          ),
          labelStyle: TextStyle(color: AppColors.foregroundDark),
        ),
      ),
      routes: {
        '/': (context) => const HomePage(),
        '/login': (context) => const LoginPage(),
      },
    );
  }
}

class HomePage extends StatelessWidget {
  const HomePage({super.key});

  void _goToLogin(BuildContext context) {
    Navigator.pushNamed(context, '/login');
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    return Scaffold(
      body: SafeArea(
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 24),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    'AREA',
                    style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                          letterSpacing: 4,
                          fontWeight: FontWeight.w600,
                        ),
                  ),
                  OutlinedButton(
                    style: OutlinedButton.styleFrom(
                      foregroundColor: colorScheme.secondary,
                      side: BorderSide(color: colorScheme.secondary),
                      shape: const StadiumBorder(),
                      padding: const EdgeInsets.symmetric(
                        horizontal: 24,
                        vertical: 12,
                      ),
                      textStyle: const TextStyle(
                        fontWeight: FontWeight.w600,
                        letterSpacing: 1.5,
                      ),
                    ),
                    onPressed: () => _goToLogin(context),
                    child: const Text('Login'),
                  ),
                ],
              ),
            ),
            Expanded(
              child: Center(
                child: ElevatedButton(
                  onPressed: () {},
                  style: ElevatedButton.styleFrom(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 48,
                      vertical: 20,
                    ),
                    shape: const StadiumBorder(),
                    textStyle: const TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.w600,
                      letterSpacing: 1.5,
                    ),
                    backgroundColor: colorScheme.primary,
                    foregroundColor: colorScheme.onPrimary,
                    elevation: 24,
                    shadowColor: Colors.black.withOpacity(0.35),
                  ),
                  child: const Text('test backend'),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  String? _status;

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  void _submit() {
    setState(() {
      _status = 'Tentative de connexion pour ${_emailController.text.trim()}';
    });
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    return Scaffold(
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 420),
              child: Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 24, vertical: 32),
                decoration: BoxDecoration(
                  color: colorScheme.surface,
                  borderRadius: BorderRadius.circular(24),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withOpacity(0.08),
                      blurRadius: 24,
                      offset: const Offset(0, 16),
                    ),
                  ],
                ),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Connexion',
                      style:
                          Theme.of(context).textTheme.headlineSmall?.copyWith(
                                fontWeight: FontWeight.bold,
                              ),
                    ),
                    const SizedBox(height: 8),
                    Text(
                      'Entrez vos identifiants pour accéder à Area',
                      style: Theme.of(context)
                          .textTheme
                          .bodyMedium
                          ?.copyWith(color: mutedColor(context)),
                    ),
                    const SizedBox(height: 24),
                    TextField(
                      controller: _emailController,
                      keyboardType: TextInputType.emailAddress,
                      decoration: const InputDecoration(
                        labelText: 'Email',
                      ),
                    ),
                    const SizedBox(height: 16),
                    TextField(
                      controller: _passwordController,
                      obscureText: true,
                      decoration: const InputDecoration(
                        labelText: 'Mot de passe',
                      ),
                    ),
                    const SizedBox(height: 24),
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton(
                        onPressed: _submit,
                        style: ElevatedButton.styleFrom(
                          backgroundColor: colorScheme.primary,
                          foregroundColor: colorScheme.onPrimary,
                          padding: const EdgeInsets.symmetric(
                            horizontal: 24,
                            vertical: 16,
                          ),
                          shape: const StadiumBorder(),
                          textStyle: const TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                        child: const Text('Se connecter'),
                      ),
                    ),
                    if (_status != null) ...[
                      const SizedBox(height: 16),
                      Container(
                        width: double.infinity,
                        padding: const EdgeInsets.symmetric(
                          horizontal: 16,
                          vertical: 12,
                        ),
                        decoration: BoxDecoration(
                          color: surfaceBorderColor(context),
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: Text(
                          _status!,
                          style: Theme.of(context)
                              .textTheme
                              .bodySmall
                              ?.copyWith(color: colorScheme.onSurface),
                        ),
                      ),
                    ],
                    const SizedBox(height: 16),
                    Align(
                      alignment: Alignment.center,
                      child: TextButton(
                        onPressed: () => Navigator.popUntil(
                          context,
                          (route) => route.settings.name == '/',
                        ),
                        child: const Text('Retour à l’accueil'),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}
