import 'package:flutter/material.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:provider/provider.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:app_links/app_links.dart';
import '../../providers/auth_provider.dart';
import '../../services/service_connector.dart';
import '../../models/service_model.dart';
import '../../theme/colors.dart';

class ServicesScreen extends StatefulWidget {
  const ServicesScreen({super.key});

  @override
  State<ServicesScreen> createState() => _ServicesScreenState();
}

class _ServicesScreenState extends State<ServicesScreen> {
  late final ServiceConnector _connector;
  final AppLinks _appLinks = AppLinks();
  final TextEditingController _searchController = TextEditingController();

  List<ServiceModel> _services = [];
  List<ServiceModel> _filteredServices = [];
  bool _isLoading = true;
  String? _error;
  String? _connectingService;

  @override
  void initState() {
    super.initState();
    _connector = ServiceConnector(
      baseUrl: dotenv.env['BASE_URL'] ?? 'http://10.0.2.2:8083',
    );
    _loadServices();
    _listenToDeepLinks();
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _listenToDeepLinks() {
    _appLinks.uriLinkStream.listen((uri) {
      if (uri.scheme == 'area' && uri.host == 'auth') {
        // OAuth callback received - reload services
        _loadServices();
        setState(() {
          _connectingService = null;
        });
      }
    });
  }

  Future<void> _loadServices() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      final authProvider = context.read<AuthProvider>();
      final user = authProvider.user;

      if (user == null || user['id'] == null) {
        throw Exception('User not authenticated');
      }

      final services = await _connector.fetchServices(user['id']);
      setState(() {
        _services = services;
        _filteredServices = services;
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
    }
  }

  void _filterServices(String query) {
    setState(() {
      if (query.isEmpty) {
        _filteredServices = _services;
      } else {
        _filteredServices = _services
            .where((s) =>
                s.displayName.toLowerCase().contains(query.toLowerCase()))
            .toList();
      }
    });
  }

  Future<void> _connectService(String serviceName) async {
    setState(() {
      _connectingService = serviceName;
    });

    try {
      final authProvider = context.read<AuthProvider>();
      final user = authProvider.user;

      if (user == null || user['id'] == null) {
        throw Exception('User not authenticated');
      }

      final authUrl =
          await _connector.getAuthUrl(serviceName, user['id']);

      if (!await launchUrl(
        Uri.parse(authUrl),
        mode: LaunchMode.externalApplication,
      )) {
        throw Exception('Could not launch browser');
      }

      // Don't reset _connectingService here - wait for deep link callback
    } catch (e) {
      setState(() {
        _connectingService = null;
      });
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Erreur: ${e.toString()}'),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  Future<void> _disconnectService(String serviceName) async {
    final confirmed = await _showDisconnectDialog(serviceName);
    if (confirmed != true) return;

    try {
      final authProvider = context.read<AuthProvider>();
      final user = authProvider.user;
      final token = authProvider.token;

      if (user == null || token == null) {
        throw Exception('User not authenticated');
      }

      await _connector.disconnectService(
          serviceName, user['id'], token);

      await _loadServices();

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Service déconnecté avec succès'),
            backgroundColor: Colors.green,
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Erreur: ${e.toString()}'),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  Future<bool?> _showDisconnectDialog(String serviceName) {
    return showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Déconnecter le service'),
        content: Text(
          'Êtes-vous sûr de vouloir déconnecter $serviceName ?',
        ),
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
            child: const Text('Déconnecter'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final connectedServices =
        _filteredServices.where((s) => s.isConnected).toList();
    final availableServices =
        _filteredServices.where((s) => !s.isConnected).toList();

    return Scaffold(
      backgroundColor: theme.colorScheme.background,
      body: SafeArea(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Header
            Padding(
              padding: const EdgeInsets.all(24.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'Services',
                    style: theme.textTheme.displayLarge,
                  ),
                  const SizedBox(height: 8),
                  Text(
                    '${connectedServices.length} connectés sur ${_services.length}',
                    style: theme.textTheme.bodyMedium?.copyWith(
                      color: AppColors.darkGrey,
                    ),
                  ),
                ],
              ),
            ),

            // Search bar
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 24.0),
              child: TextField(
                controller: _searchController,
                onChanged: _filterServices,
                decoration: InputDecoration(
                  hintText: 'Rechercher un service...',
                  prefixIcon: const Icon(Icons.search),
                  suffixIcon: _searchController.text.isNotEmpty
                      ? IconButton(
                          icon: const Icon(Icons.clear),
                          onPressed: () {
                            _searchController.clear();
                            _filterServices('');
                          },
                        )
                      : null,
                ),
              ),
            ),

            const SizedBox(height: 24),

            // Content
            Expanded(
              child: _isLoading
                  ? const Center(child: CircularProgressIndicator())
                  : _error != null
                      ? Center(
                          child: Column(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              const Icon(
                                Icons.error_outline,
                                size: 64,
                                color: AppColors.darkGrey,
                              ),
                              const SizedBox(height: 16),
                              Text(
                                'Erreur de chargement',
                                style: theme.textTheme.titleMedium,
                              ),
                              const SizedBox(height: 8),
                              Text(
                                _error!,
                                style: theme.textTheme.bodySmall,
                                textAlign: TextAlign.center,
                              ),
                              const SizedBox(height: 24),
                              ElevatedButton.icon(
                                onPressed: _loadServices,
                                icon: const Icon(Icons.refresh),
                                label: const Text('Réessayer'),
                              ),
                            ],
                          ),
                        )
                      : RefreshIndicator(
                          onRefresh: _loadServices,
                          child: ListView(
                            padding: const EdgeInsets.symmetric(
                                horizontal: 24),
                            children: [
                              // Connected services
                              if (connectedServices.isNotEmpty) ...[
                                Text(
                                  'Services connectés',
                                  style: theme.textTheme.titleMedium,
                                ),
                                const SizedBox(height: 12),
                                ...connectedServices.map(
                                  (service) => _ServiceCard(
                                    service: service,
                                    isConnecting: _connectingService ==
                                        service.name,
                                    onConnect: () =>
                                        _connectService(service.name),
                                    onDisconnect: () =>
                                        _disconnectService(service.name),
                                  ),
                                ),
                                const SizedBox(height: 32),
                              ],

                              // Available services
                              if (availableServices.isNotEmpty) ...[
                                Text(
                                  'Services disponibles',
                                  style: theme.textTheme.titleMedium,
                                ),
                                const SizedBox(height: 12),
                                ...availableServices.map(
                                  (service) => _ServiceCard(
                                    service: service,
                                    isConnecting: _connectingService ==
                                        service.name,
                                    onConnect: () =>
                                        _connectService(service.name),
                                    onDisconnect: () =>
                                        _disconnectService(service.name),
                                  ),
                                ),
                              ],

                              if (_filteredServices.isEmpty) ...[
                                const SizedBox(height: 64),
                                Center(
                                  child: Column(
                                    children: [
                                      const Icon(
                                        Icons.search_off,
                                        size: 64,
                                        color: AppColors.darkGrey,
                                      ),
                                      const SizedBox(height: 16),
                                      Text(
                                        'Aucun service trouvé',
                                        style: theme.textTheme.titleMedium,
                                      ),
                                    ],
                                  ),
                                ),
                              ],

                              const SizedBox(height: 24),
                            ],
                          ),
                        ),
            ),
          ],
        ),
      ),
    );
  }
}

class _ServiceCard extends StatelessWidget {
  final ServiceModel service;
  final bool isConnecting;
  final VoidCallback onConnect;
  final VoidCallback onDisconnect;

  const _ServiceCard({
    required this.service,
    required this.isConnecting,
    required this.onConnect,
    required this.onDisconnect,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Row(
          children: [
            // Service icon
            Container(
              width: 48,
              height: 48,
              decoration: BoxDecoration(
                color: AppColors.lightGrey,
                borderRadius: BorderRadius.circular(12),
              ),
              child: Center(
                child: _getServiceIcon(service.name),
              ),
            ),

            const SizedBox(width: 16),

            // Service name and status
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    service.displayName,
                    style: theme.textTheme.titleMedium?.copyWith(
                      fontSize: 16,
                    ),
                  ),
                  if (service.isConnected) ...[
                    const SizedBox(height: 4),
                    Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 8,
                        vertical: 2,
                      ),
                      decoration: BoxDecoration(
                        color: Colors.green.shade50,
                        borderRadius: BorderRadius.circular(4),
                        border: Border.all(
                          color: Colors.green.shade200,
                        ),
                      ),
                      child: Text(
                        'Connecté',
                        style: TextStyle(
                          fontSize: 12,
                          color: Colors.green.shade700,
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                    ),
                  ],
                ],
              ),
            ),

            const SizedBox(width: 12),

            // Connect/Disconnect button
            if (isConnecting)
              const SizedBox(
                width: 24,
                height: 24,
                child: CircularProgressIndicator(strokeWidth: 2),
              )
            else if (service.isConnected)
              OutlinedButton(
                onPressed: onDisconnect,
                style: OutlinedButton.styleFrom(
                  foregroundColor: Colors.red.shade600,
                  side: BorderSide(color: Colors.red.shade600),
                  padding: const EdgeInsets.symmetric(
                    horizontal: 16,
                    vertical: 8,
                  ),
                ),
                child: const Text('Déconnecter'),
              )
            else
              ElevatedButton(
                onPressed: onConnect,
                style: ElevatedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 16,
                    vertical: 8,
                  ),
                ),
                child: const Text('Connecter'),
              ),
          ],
        ),
      ),
    );
  }

  Widget _getServiceIcon(String serviceName) {
    final iconMap = {
      'google': Icons.g_mobiledata,
      'github': Icons.code,
      'discord': Icons.chat,
      'spotify': Icons.music_note,
      'notion': Icons.note,
      'slack': Icons.tag,
      'twitter': Icons.tag,
      'gmail': Icons.email,
    };

    final colorMap = {
      'google': Colors.red,
      'github': Colors.black,
      'discord': Color(0xFF5865F2),
      'spotify': Color(0xFF1DB954),
      'notion': Colors.black,
      'slack': Color(0xFF4A154B),
      'twitter': Color(0xFF1DA1F2),
      'gmail': Colors.red,
    };

    return Icon(
      iconMap[serviceName.toLowerCase()] ?? Icons.widgets,
      color: colorMap[serviceName.toLowerCase()] ?? AppColors.deepBlue,
      size: 28,
    );
  }
}