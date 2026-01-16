import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:provider/provider.dart';

import '../../providers/auth_provider.dart';
import '../../providers/theme_provider.dart';
import '../../services/config_service.dart';
import '../../theme/colors.dart';
import 'support_screen.dart';

class ProfileScreen extends StatefulWidget {
  const ProfileScreen({super.key});

  @override
  State<ProfileScreen> createState() => _ProfileScreenState();
}

class _ProfileScreenState extends State<ProfileScreen> {
  final FlutterSecureStorage _storage = const FlutterSecureStorage();

  final TextEditingController _ipController = TextEditingController();
  final TextEditingController _nameController = TextEditingController();
  final TextEditingController _usernameController = TextEditingController();
  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();
  final TextEditingController _avatarUrlController = TextEditingController();

  bool _isSavingServer = false;
  bool _notifyApp = true;
  bool _notifyEmail = true;
  List<_ProfileNotification> _notifications = [];

  static const _keyProfileDraft = 'profile_draft';
  static const _keyNotifications = 'profile_notifications';
  static const _keyNotifyApp = 'profile_notify_app';
  static const _keyNotifyEmail = 'profile_notify_email';

  @override
  void initState() {
    super.initState();
    _loadServerIp();
    _hydrateProfile();
  }

  Future<void> _loadServerIp() async {
    final ip = await ConfigService.getServerIp();
    if (mounted) {
      setState(() => _ipController.text = ip ?? '');
    }
  }

  Future<void> _saveServerIp() async {
    setState(() => _isSavingServer = true);
    await ConfigService.setServerIp(_ipController.text);
    setState(() => _isSavingServer = false);

    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Server IP sauvegardée avec succès'),
          backgroundColor: Colors.green,
        ),
      );
    }
  }

  Future<void> _hydrateProfile() async {
    // Profile draft
    final rawDraft = await _storage.read(key: _keyProfileDraft);
    if (rawDraft != null) {
      try {
        final decoded = jsonDecode(rawDraft) as Map<String, dynamic>;
        _nameController.text = decoded['name'] ?? '';
        _usernameController.text = decoded['username'] ?? '';
        _emailController.text = decoded['email'] ?? '';
        _passwordController.text = decoded['password'] ?? '';
        _avatarUrlController.text = decoded['avatarUrl'] ?? '';
      } catch (_) {
        // ignore invalid data
      }
    }

    // Notifications
    final rawNotif = await _storage.read(key: _keyNotifications);
    if (rawNotif != null) {
      try {
        final decoded = jsonDecode(rawNotif) as List<dynamic>;
        _notifications = decoded
            .map((e) => _ProfileNotification.fromJson(
                Map<String, dynamic>.from(e as Map)))
            .toList();
      } catch (_) {
        _notifications = _seedNotifications();
      }
    } else {
      _notifications = _seedNotifications();
    }

    // Toggles
    final notifyApp = await _storage.read(key: _keyNotifyApp);
    final notifyEmail = await _storage.read(key: _keyNotifyEmail);
    _notifyApp = notifyApp == null ? true : notifyApp == 'true';
    _notifyEmail = notifyEmail == null ? true : notifyEmail == 'true';

    if (mounted) setState(() {});
  }

  List<_ProfileNotification> _seedNotifications() {
    final seed = [
      _ProfileNotification(
        id: 'seed-1',
        title: 'Bienvenue sur AREA',
        detail:
            'Connectez un service pour commencer à créer des automatisations.',
        type: NotificationType.info,
        createdAt: DateTime.now(),
      ),
    ];
    _persistNotifications(seed);
    return seed;
  }

  Future<void> _persistDraft() async {
    final draft = {
      'name': _nameController.text,
      'username': _usernameController.text,
      'email': _emailController.text,
      'password': _passwordController.text,
      'avatarUrl': _avatarUrlController.text,
    };
    await _storage.write(key: _keyProfileDraft, value: jsonEncode(draft));
  }

  Future<void> _persistNotifications(
      List<_ProfileNotification> notifications) async {
    await _storage.write(
      key: _keyNotifications,
      value: jsonEncode(notifications.map((n) => n.toJson()).toList()),
    );
  }

  Future<void> _persistToggles() async {
    await _storage.write(key: _keyNotifyApp, value: _notifyApp.toString());
    await _storage.write(key: _keyNotifyEmail, value: _notifyEmail.toString());
  }

  void _openEditSheet(Map<String, dynamic>? user) {
    final colors = context.appColors;
    final displayName =
        _nameController.text.isNotEmpty ? _nameController.text : _displayName(user);
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) {
        return DraggableScrollableSheet(
          initialChildSize: 0.82,
          minChildSize: 0.6,
          maxChildSize: 0.95,
          builder: (context, scrollController) {
            return Container(
              decoration: BoxDecoration(
                color: colors.white,
                borderRadius:
                    const BorderRadius.vertical(top: Radius.circular(20)),
                boxShadow: [
                  BoxShadow(
                    color: colors.grey.withOpacity(0.25),
                    blurRadius: 18,
                    offset: const Offset(0, -4),
                  ),
                ],
              ),
              child: ListView(
                controller: scrollController,
                padding:
                    const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
                children: [
                  Center(
                    child: Container(
                      width: 46,
                      height: 4,
                      decoration: BoxDecoration(
                        color: colors.grey,
                        borderRadius: BorderRadius.circular(8),
                      ),
                    ),
                  ),
                  const SizedBox(height: 12),
                  Text(
                    'Modifier le profil',
                    style: Theme.of(context).textTheme.titleLarge,
                  ),
                  Text(
                    'Mises à jour locales uniquement (pas d’appel backend).',
                    style: Theme.of(context)
                        .textTheme
                        .bodySmall
                        ?.copyWith(color: colors.darkGrey),
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      Container(
                        width: 64,
                        height: 64,
                        decoration: BoxDecoration(
                          shape: BoxShape.circle,
                          color: colors.deepBlue.withOpacity(0.12),
                          border: Border.all(color: colors.deepBlue),
                        ),
                        child: _avatarUrlController.text.isNotEmpty
                            ? ClipOval(
                                child: Image.network(
                                  _avatarUrlController.text,
                                  fit: BoxFit.cover,
                                  errorBuilder: (_, __, ___) => _initialsAvatar(displayName),
                                ),
                              )
                            : _initialsAvatar(displayName),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: ElevatedButton.icon(
                          onPressed: () {
                            // Pour un flux réel, brancher un file picker ici (ex: file_picker).
                            _addNotification(
                              title: 'Fichier avatar',
                              detail:
                                  'Sélection de fichier simulée (brancher file_picker pour un vrai flux).',
                              type: NotificationType.info,
                            );
                          },
                          icon: const Icon(Icons.upload_file_outlined),
                          label: const Text('Choisir un fichier'),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  TextFormField(
                    controller: _nameController,
                    decoration: const InputDecoration(
                      labelText: 'Nom / affichage',
                      prefixIcon: Icon(Icons.person_outline),
                    ),
                  ),
                  const SizedBox(height: 12),
                  TextFormField(
                    controller: _usernameController,
                    decoration: const InputDecoration(
                      labelText: 'Username',
                      prefixIcon: Icon(Icons.badge_outlined),
                    ),
                  ),
                  const SizedBox(height: 12),
                  TextFormField(
                    controller: _emailController,
                    decoration: const InputDecoration(
                      labelText: 'Email',
                      prefixIcon: Icon(Icons.email_outlined),
                    ),
                  ),
                  const SizedBox(height: 12),
                  TextFormField(
                    controller: _passwordController,
                    obscureText: true,
                    decoration: const InputDecoration(
                      labelText: 'Mot de passe (stocké localement)',
                      prefixIcon: Icon(Icons.lock_outline),
                    ),
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      Expanded(
                        child: OutlinedButton(
                          onPressed: () {
                            Navigator.of(context).pop();
                          },
                          child: const Text('Annuler'),
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: ElevatedButton(
                          onPressed: () {
                            _persistDraft();
                            _addNotification(
                              title: 'Profil mis à jour',
                              detail:
                                  'Modifications enregistrées localement (sans backend).',
                              type: NotificationType.info,
                            );
                            Navigator.of(context).pop();
                          },
                          child: const Text('Enregistrer'),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            );
          },
        );
      },
    ).whenComplete(() => setState(() {}));
  }

  Widget _initialsAvatar(String displayName) {
    return Center(
      child: Text(
        displayName.isNotEmpty
            ? displayName.substring(0, displayName.length >= 2 ? 2 : 1)
            : 'A',
        style: const TextStyle(
          fontWeight: FontWeight.w700,
          color: Colors.white,
        ),
      ),
    );
  }

  void _addNotification({
    required String title,
    required String detail,
    required NotificationType type,
  }) {
    final next = _ProfileNotification(
      id: 'notif-${DateTime.now().millisecondsSinceEpoch}',
      title: title,
      detail: detail,
      type: type,
      createdAt: DateTime.now(),
    );
    setState(() {
      _notifications = [next, ..._notifications].take(25).toList();
    });
    _persistNotifications(_notifications);
  }

  void _clearNotifications() {
    setState(() {
      _notifications = [];
    });
    _persistNotifications(_notifications);
  }

  String _maskToken(String? token) {
    if (token == null || token.isEmpty) return 'Aucun token actif';
    if (token.length <= 10) return token;
    return '${token.substring(0, 8)}…${token.substring(token.length - 4)}';
  }

  @override
  void dispose() {
    _ipController.dispose();
    _nameController.dispose();
    _usernameController.dispose();
    _emailController.dispose();
    _passwordController.dispose();
    _avatarUrlController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colors = context.appColors;
    final authProvider = context.watch<AuthProvider>();
    final themeProvider = context.watch<ThemeProvider>();
    final user = authProvider.user;
    final token = authProvider.token;

    final displayName = _nameController.text.isNotEmpty
        ? _nameController.text
        : _displayName(user);
    final displayEmail =
        _emailController.text.isNotEmpty ? _emailController.text : _email(user);
    final username = _usernameController.text.isNotEmpty
        ? _usernameController.text
        : user?['username']?.toString() ?? 'N/A';
    final avatarUrl = _avatarUrlController.text.isNotEmpty
        ? _avatarUrlController.text
        : user?['avatarUrl']?.toString() ?? '';

    final isDark = theme.brightness == Brightness.dark;
    final surface = isDark ? colors.white.withOpacity(0.08) : Colors.white;
    final canvas = colors.white;

    return Scaffold(
      backgroundColor: canvas,
      appBar: AppBar(
        title: const Text('Profil'),
        backgroundColor: Colors.transparent,
        elevation: 0,
        foregroundColor: isDark ? Colors.white : colors.almostBlack,
        actions: [
          IconButton(
            tooltip: 'Déconnexion',
            icon: const Icon(Icons.power_settings_new_rounded),
            onPressed: () => _showLogoutDialog(context),
          ),
        ],
      ),
      body: Stack(
        children: [
          SafeArea(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(20.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Center(
                    child: _HeroCard(
                      displayName: displayName,
                      displayEmail: displayEmail,
                      avatarUrl: avatarUrl,
                      token: token,
                      isDark: isDark,
                      colors: colors,
                      fallbackAvatar: _initialsAvatar(displayName),
                    ),
                  ),
                  const SizedBox(height: 20),

                  _SectionTile(
                    title: 'Compte',
                    initiallyExpanded: true,
                    surface: surface,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        _SettingTile(
                          icon: Icons.manage_accounts_outlined,
                          title: 'Modifier le profil',
                          subtitle: displayName,
                          onTap: () => _openEditSheet(user),
                          isDark: isDark,
                        ),
                        const Divider(height: 1),
                        _SettingTile(
                          icon: Icons.lock_outline,
                          title: 'Mot de passe & sécurité',
                          subtitle: 'Non synchronisé (local)',
                          onTap: () => _addNotification(
                            title: 'Sécurité',
                            detail:
                                'Aucune route backend : mise à jour simulée uniquement.',
                            type: NotificationType.info,
                          ),
                          isDark: isDark,
                        ),
                        const Divider(height: 1),
                        _SettingTile(
                          icon: Icons.language_rounded,
                          title: 'Langue',
                          trailingText: 'FR',
                          onTap: () {},
                          isDark: isDark,
                        ),
                      ],
                    ),
                  ),

                  const SizedBox(height: 12),

                  _SectionTile(
                    title: 'Notifications',
                    initiallyExpanded: true,
                    surface: surface,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        OutlinedButton.icon(
                          onPressed: () {
                            _addNotification(
                              title: 'AREA créée',
                              detail: 'Votre nouvelle automation est prête.',
                              type: NotificationType.areaCreated,
                            );
                          },
                          icon: const Icon(Icons.auto_awesome_outlined),
                          label: const Text('Simuler création AREA'),
                          style: OutlinedButton.styleFrom(
                            foregroundColor:
                                isDark ? Colors.white : colors.deepBlue,
                            side: BorderSide(
                              color: isDark
                                  ? Colors.white24
                                  : colors.grey.withOpacity(0.6),
                            ),
                          ),
                        ),
                        const SizedBox(height: 12),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(
                              'Feed',
                              style: theme.textTheme.bodyMedium,
                            ),
                            TextButton(
                              onPressed:
                                  _notifications.isEmpty ? null : _clearNotifications,
                              child: const Text('Tout effacer'),
                            ),
                          ],
                        ),
                        const SizedBox(height: 8),
                        if (_notifications.isEmpty)
                          Text(
                            'Aucune notification locale.',
                            style: theme.textTheme.bodySmall
                                ?.copyWith(color: colors.darkGrey),
                          )
                        else
                          Column(
                          children: _notifications
                              .map((n) => _NotificationTile(notification: n))
                              .toList(),
                        ),
                      ],
                    ),
                  ),

                  const SizedBox(height: 12),

                  _SectionTile(
                    title: 'Préférences',
                    surface: surface,
                    child: Column(
                      children: [
                        _SettingTile(
                          icon: Icons.info_outline,
                          title: 'À propos',
                          subtitle: 'AREA mobile',
                          onTap: () => _addNotification(
                            title: 'À propos',
                            detail: 'Version 1.0.0 — build local.',
                            type: NotificationType.info,
                          ),
                          isDark: isDark,
                        ),
                        const Divider(height: 1),
                        _SettingTile(
                          icon: Icons.dark_mode_outlined,
                          title: 'Thème',
                          subtitle: isDark ? 'Sombre (système)' : 'Clair (système)',
                          onTap: () {},
                          isDark: isDark,
                        ),
                        const Divider(height: 1),
                        _SettingTile(
                          icon: Icons.event_available_outlined,
                          title: 'Rendez-vous',
                          subtitle: 'Bientôt disponible',
                          onTap: () {},
                          isDark: isDark,
                        ),
                      ],
                    ),
                  ),

                  const SizedBox(height: 12),

                  _SectionTile(
                    title: 'Affichage & accessibilité',
                    surface: surface,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Correction des couleurs',
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: colors.darkGrey,
                          ),
                        ),
                        const SizedBox(height: 8),
                        Wrap(
                          spacing: 8,
                          runSpacing: 8,
                          children: [
                            ChoiceChip(
                              label: const Text('Standard'),
                              selected:
                                  themeProvider.visionMode == VisionMode.normal,
                              onSelected: (v) {
                                if (v) {
                                  themeProvider
                                      .setVisionMode(VisionMode.normal);
                                }
                              },
                              selectedColor: colors.deepBlue,
                              backgroundColor: colors.lightGrey,
                              labelStyle: TextStyle(
                                color: themeProvider.visionMode ==
                                        VisionMode.normal
                                    ? Colors.white
                                    : colors.almostBlack,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                            ChoiceChip(
                              label: const Text('Tritanopie'),
                              selected: themeProvider.visionMode ==
                                  VisionMode.tritanopia,
                              onSelected: (v) {
                                if (v) {
                                  themeProvider
                                      .setVisionMode(VisionMode.tritanopia);
                                }
                              },
                              selectedColor: colors.deepBlue,
                              backgroundColor: colors.lightGrey,
                              labelStyle: TextStyle(
                                color: themeProvider.visionMode ==
                                        VisionMode.tritanopia
                                    ? Colors.white
                                    : colors.almostBlack,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 12),
                        Row(
                          children: [
                            Icon(
                              theme.brightness == Brightness.dark
                                  ? Icons.dark_mode_outlined
                                  : Icons.light_mode_outlined,
                              color: colors.deepBlue,
                            ),
                            const SizedBox(width: 8),
                            Expanded(
                              child: Text(
                                'Le thème suit le réglage système (clair/sombre).',
                                style: theme.textTheme.bodySmall?.copyWith(
                                  color: colors.darkGrey,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),

                  const SizedBox(height: 12),

                  _SectionTile(
                    title: 'Configuration serveur',
                    surface: surface,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        TextFormField(
                          controller: _ipController,
                          decoration: const InputDecoration(
                            labelText: 'Adresse du serveur (IP)',
                            prefixIcon: Icon(Icons.cloud_outlined),
                          ),
                        ),
                        const SizedBox(height: 12),
                        SizedBox(
                          width: double.infinity,
                          height: 45,
                          child: ElevatedButton.icon(
                            onPressed: _isSavingServer ? null : _saveServerIp,
                            icon: _isSavingServer
                                ? const SizedBox(
                                    width: 16,
                                    height: 16,
                                    child: CircularProgressIndicator(
                                      strokeWidth: 2,
                                      color: Colors.white,
                                    ),
                                  )
                                : const Icon(Icons.save),
                            label: const Text('Sauvegarder'),
                          ),
                        ),
                      ],
                    ),
                  ),

                  const SizedBox(height: 12),

                  _SectionTile(
                    title: 'Préférences de notification',
                    surface: surface,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        SwitchListTile(
                          title: const Text('Notifications in-app'),
                          value: _notifyApp,
                          onChanged: (v) {
                            setState(() => _notifyApp = v);
                            _persistToggles();
                          },
                        ),
                        SwitchListTile(
                          title: const Text('Notifications email'),
                          value: _notifyEmail,
                          onChanged: (v) {
                            setState(() => _notifyEmail = v);
                            _persistToggles();
                          },
                        ),
                      ],
                    ),
                  ),

                  const SizedBox(height: 12),

                  _SectionTile(
                    title: 'Support',
                    surface: surface,
                    child: Column(
                      children: [
                        ListTile(
                          contentPadding: const EdgeInsets.symmetric(horizontal: 4),
                          leading: Icon(Icons.help_outline, color: isDark ? Colors.white : colors.deepBlue),
                          title: const Text('Help Center'),
                          trailing: const Icon(Icons.chevron_right),
                          onTap: () {
                            Navigator.of(context).push(
                              MaterialPageRoute(builder: (_) => const SupportScreen()),
                            );
                          },
                        ),
                      ],
                    ),
                  ),

                  const SizedBox(height: 24),
                  // Actions bas de page
                  Column(
                    children: [
                      Row(
                        children: [
                          Expanded(
                            child: OutlinedButton.icon(
                              onPressed: () {
                                _addNotification(
                                  title: 'Suppression (simulée)',
                                  detail:
                                      'Aucune action backend, suppression non effectuée.',
                                  type: NotificationType.warning,
                                );
                              },
                              icon: const Icon(Icons.delete_outline),
                              label: const Text('Supprimer le compte'),
                              style: OutlinedButton.styleFrom(
                                foregroundColor: Colors.red.shade600,
                                side: BorderSide(color: Colors.red.shade300),
                              ),
                            ),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: OutlinedButton.icon(
                              onPressed: () async {
                                await context.read<AuthProvider>().logout();
                                if (mounted) {
                                  Navigator.of(context).pushNamedAndRemoveUntil(
                                    '/login',
                                    (route) => false,
                                  );
                                }
                              },
                              icon: const Icon(Icons.swap_horiz),
                              label: const Text('Changer de compte'),
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 12),
                      SizedBox(
                        width: double.infinity,
                        height: 50,
                        child: ElevatedButton.icon(
                          onPressed: () async {
                            final confirmed = await _showLogoutDialog(context);
                            if (confirmed == true && context.mounted) {
                              await context.read<AuthProvider>().logout();
                              if (context.mounted) {
                                Navigator.of(context).pushNamedAndRemoveUntil(
                                  '/login',
                                  (route) => false,
                                );
                              }
                            }
                          },
                          icon: const Icon(Icons.logout),
                          label: const Text('Se déconnecter'),
                          style: ElevatedButton.styleFrom(
                            backgroundColor: colors.deepBlue,
                            foregroundColor: Colors.white,
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(10),
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),

                  const SizedBox(height: 16),
                  Center(
                    child: Text(
                      'Version 1.0.0',
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: colors.darkGrey,
                      ),
                    ),
                  ),
                  const SizedBox(height: 8),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Future<bool?> _showLogoutDialog(BuildContext context) {
    return showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Déconnexion'),
        content: const Text('Êtes-vous sûr de vouloir vous déconnecter ?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('Annuler'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.of(context).pop(true),
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.red.shade600,
              foregroundColor: Colors.white,
            ),
            child: const Text('Se déconnecter'),
          ),
        ],
      ),
    );
  }

  String _displayName(Map<String, dynamic>? user) {
    return user?['name']?.toString().isNotEmpty == true
        ? user!['name'].toString()
        : user?['username']?.toString().isNotEmpty == true
            ? user!['username'].toString()
            : user?['email']?.toString() ?? 'Utilisateur';
  }

  String _email(Map<String, dynamic>? user) {
    return user?['email']?.toString() ?? 'Email indisponible';
  }
}

class _SectionTile extends StatelessWidget {
  final String title;
  final Widget child;
  final bool initiallyExpanded;
  final Color? surface;

  const _SectionTile({
    required this.title,
    required this.child,
    this.initiallyExpanded = false,
    this.surface,
  });

  @override
  Widget build(BuildContext context) {
    final colors = context.appColors;
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final borderColor = isDark
        ? Colors.white.withOpacity(0.12)
        : colors.deepBlue.withOpacity(0.12);
    return Container(
      decoration: BoxDecoration(
        color: surface ?? colors.white,
        borderRadius: BorderRadius.circular(14),
        boxShadow: [
          BoxShadow(
            color: colors.grey.withOpacity(0.12),
            blurRadius: 14,
            offset: const Offset(0, 6),
          ),
        ],
        border: Border.all(color: borderColor),
      ),
      child: Theme(
        data: Theme.of(context).copyWith(
          dividerColor: Colors.transparent,
          listTileTheme: ListTileThemeData(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(14),
            ),
            iconColor: isDark ? Colors.white70 : colors.almostBlack,
            textColor: isDark ? Colors.white : colors.almostBlack,
          ),
        ),
        child: ExpansionTile(
          initiallyExpanded: initiallyExpanded,
          tilePadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
          childrenPadding: const EdgeInsets.symmetric(
            horizontal: 16,
            vertical: 10,
          ),
          title: Text(
            title,
            style: Theme.of(context)
                .textTheme
                .titleMedium
                ?.copyWith(color: isDark ? Colors.white : null),
          ),
          backgroundColor: surface ?? colors.white,
          collapsedBackgroundColor: surface ?? colors.white,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(14),
          ),
          collapsedShape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(14),
          ),
          children: [child],
        ),
      ),
    );
  }
}

class _InfoRow extends StatelessWidget {
  final IconData icon;
  final String label;
  final String value;

  const _InfoRow({
    required this.icon,
    required this.label,
    required this.value,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colors = context.appColors;
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Row(
      children: [
        Icon(icon, color: isDark ? Colors.white : colors.deepBlue, size: 24),
        const SizedBox(width: 16),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                label,
                style: theme.textTheme.bodySmall?.copyWith(
                  color: isDark ? Colors.white70 : colors.darkGrey,
                ),
              ),
              const SizedBox(height: 4),
              Text(
                value,
                style: theme.textTheme.bodyLarge?.copyWith(
                  fontWeight: FontWeight.w500,
                  color: isDark ? Colors.white : null,
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }
}

class _StatusPill extends StatelessWidget {
  final String label;
  final Color color;

  const _StatusPill({required this.label, required this.color});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: color.withOpacity(0.12),
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: color.withOpacity(0.5)),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: 8,
            height: 8,
            decoration: BoxDecoration(
              color: color,
              shape: BoxShape.circle,
            ),
          ),
          const SizedBox(width: 6),
          Text(
            label,
            style: Theme.of(context)
                .textTheme
                .bodySmall
                ?.copyWith(color: color, fontWeight: FontWeight.w600),
          ),
        ],
      ),
    );
  }
}

class _HeroCard extends StatelessWidget {
  final String displayName;
  final String displayEmail;
  final String avatarUrl;
  final String? token;
  final bool isDark;
  final AppColorPalette colors;
  final Widget fallbackAvatar;

  const _HeroCard({
    required this.displayName,
    required this.displayEmail,
    required this.avatarUrl,
    required this.token,
    required this.isDark,
    required this.colors,
    required this.fallbackAvatar,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Container(
          width: 110,
          height: 110,
          decoration: BoxDecoration(
            shape: BoxShape.circle,
            color: isDark
                ? Colors.white.withOpacity(0.08)
                : colors.deepBlue.withOpacity(0.08),
            border: Border.all(
              color:
                  isDark ? Colors.white.withOpacity(0.25) : colors.deepBlue,
              width: 1.6,
            ),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(0.20),
                blurRadius: 16,
                offset: const Offset(0, 8),
              ),
            ],
          ),
          child: avatarUrl.isNotEmpty
              ? ClipOval(
                  child: Image.network(
                    avatarUrl,
                    fit: BoxFit.cover,
                    errorBuilder: (_, __, ___) => fallbackAvatar,
                  ),
                )
              : fallbackAvatar,
        ),
        const SizedBox(height: 16),
        Column(
          children: [
            Text(
              displayName,
              style: Theme.of(context)
                  .textTheme
                  .titleLarge
                  ?.copyWith(color: isDark ? Colors.white : colors.almostBlack),
            ),
            const SizedBox(height: 4),
            Text(
              displayEmail,
              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: isDark ? Colors.white70 : colors.darkGrey,
                  ),
            ),
            const SizedBox(height: 8),
            _StatusPill(
              label: token != null ? 'Authentifié' : 'Hors ligne',
              color:
                  token != null ? colors.deepBlue : Colors.orange.shade700,
            ),
          ],
        ),
      ],
    );
  }
}

class _SettingTile extends StatelessWidget {
  final IconData icon;
  final String title;
  final String? subtitle;
  final String? trailingText;
  final VoidCallback onTap;
  final bool isDark;

  const _SettingTile({
    required this.icon,
    required this.title,
    required this.onTap,
    required this.isDark,
    this.subtitle,
    this.trailingText,
  });

  @override
  Widget build(BuildContext context) {
    final colors = context.appColors;
    return ListTile(
      contentPadding: const EdgeInsets.symmetric(horizontal: 4),
      leading: Icon(icon, color: isDark ? Colors.white : colors.deepBlue),
      title: Text(
        title,
        style: Theme.of(context).textTheme.bodyMedium?.copyWith(
              color: isDark ? Colors.white : colors.almostBlack,
              fontWeight: FontWeight.w600,
            ),
      ),
      subtitle: subtitle != null
          ? Text(
              subtitle!,
              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: isDark ? Colors.white70 : colors.darkGrey,
                  ),
            )
          : null,
      trailing: trailingText != null
          ? Text(
              trailingText!,
              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: isDark ? Colors.white70 : colors.darkGrey,
                  ),
            )
          : const Icon(Icons.chevron_right),
      onTap: onTap,
    );
  }
}

enum NotificationType { areaCreated, info, warning }

class _ProfileNotification {
  final String id;
  final String title;
  final String detail;
  final NotificationType type;
  final DateTime createdAt;

  _ProfileNotification({
    required this.id,
    required this.title,
    required this.detail,
    required this.type,
    required this.createdAt,
  });

  Map<String, dynamic> toJson() => {
        'id': id,
        'title': title,
        'detail': detail,
        'type': type.name,
        'createdAt': createdAt.toIso8601String(),
      };

  factory _ProfileNotification.fromJson(Map<String, dynamic> json) {
    return _ProfileNotification(
      id: json['id'] as String,
      title: json['title'] as String,
      detail: json['detail'] as String,
      type: NotificationType.values.firstWhere(
        (t) => t.name == json['type'],
        orElse: () => NotificationType.info,
      ),
      createdAt: DateTime.tryParse(json['createdAt'] as String? ?? '') ??
          DateTime.now(),
    );
  }
}

class _NotificationTile extends StatelessWidget {
  final _ProfileNotification notification;

  const _NotificationTile({required this.notification});

  @override
  Widget build(BuildContext context) {
    final colors = context.appColors;
    final isDark = Theme.of(context).brightness == Brightness.dark;

    Color tone;
    IconData icon;
    switch (notification.type) {
      case NotificationType.areaCreated:
        tone = Colors.green.shade400;
        icon = Icons.auto_awesome;
        break;
      case NotificationType.warning:
        tone = Colors.orange.shade400;
        icon = Icons.warning_amber_rounded;
        break;
      default:
        tone = isDark ? Colors.white70 : colors.deepBlue;
        icon = Icons.info_outline;
    }

    return Container(
      margin: const EdgeInsets.symmetric(vertical: 6),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: isDark ? Colors.white.withOpacity(0.06) : colors.lightGrey,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: isDark ? Colors.white12 : colors.grey,
        ),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(icon, color: tone),
          const SizedBox(width: 10),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  notification.title,
                  style: Theme.of(context)
                      .textTheme
                      .bodyMedium
                      ?.copyWith(
                          fontWeight: FontWeight.w700,
                          color: isDark ? Colors.white : null),
                ),
                const SizedBox(height: 4),
                Text(
                  notification.detail,
                  style: Theme.of(context)
                      .textTheme
                      .bodySmall
                      ?.copyWith(
                          color:
                              isDark ? Colors.white70 : colors.darkGrey),
                ),
                const SizedBox(height: 4),
                Text(
                  _formatDate(notification.createdAt),
                  style: Theme.of(context)
                      .textTheme
                      .labelSmall
                      ?.copyWith(
                          color:
                              isDark ? Colors.white60 : colors.darkGrey),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  String _formatDate(DateTime date) {
    final day = date.day.toString().padLeft(2, '0');
    final month = date.month.toString().padLeft(2, '0');
    final hour = date.hour.toString().padLeft(2, '0');
    final minute = date.minute.toString().padLeft(2, '0');
    return '$day/$month/${date.year} · $hour:$minute';
  }
}
