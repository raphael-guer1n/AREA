import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:app_links/app_links.dart';
import '../../providers/auth_provider.dart';
import '../../services/service_connector.dart';
import '../../models/service_model.dart';
import '../../theme/colors.dart';

enum SortOption { alphabetical, reverseAlphabetical }

class ServicesScreen extends StatefulWidget {
  const ServicesScreen({super.key});

  @override
  State<ServicesScreen> createState() => _ServicesScreenState();
}

class _ServicesScreenState extends State<ServicesScreen> {
  late final ServiceConnector _connector;
  final AppLinks _appLinks = AppLinks();
  final TextEditingController _searchController = TextEditingController();
  final TextEditingController _connectSearchController =
      TextEditingController();

  List<ServiceModel> _services = [];
  List<ServiceModel> _connectedServices = [];
  List<ServiceModel> _availableServices = [];
  List<ServiceModel> _visibleConnected = [];
  bool _isLoading = true;
  String? _error;
  String? _connectingService;
  bool _showFilters = false;
  SortOption _sort = SortOption.alphabetical;

  @override
  void initState() {
    super.initState();
    _connector = ServiceConnector();
    _loadServices();
    _listenToDeepLinks();
  }

  @override
  void dispose() {
    _searchController.dispose();
    _connectSearchController.dispose();
    super.dispose();
  }

  void _listenToDeepLinks() {
    _appLinks.uriLinkStream.listen((uri) {
      if (uri.scheme == 'area' && uri.host == 'auth') {
        debugPrint('[DEEP LINK] OAuth callback: $uri');
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
        _connectedServices =
            services.where((s) => s.isConnected).toList(growable: false);
        _availableServices =
            services.where((s) => !s.isConnected).toList(growable: false);
        _visibleConnected = List<ServiceModel>.from(_connectedServices);
        _applySort(_visibleConnected);
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
    final filteredConnected = _connectedServices
        .where(
          (s) => s.displayName.toLowerCase().contains(query.toLowerCase()),
        )
        .toList();
    _applySort(filteredConnected);
    setState(() {
      _visibleConnected = filteredConnected;
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

      final authUrl = await _connector.getAuthUrl(serviceName, user['id']);

      if (!await launchUrl(
        Uri.parse(authUrl),
        mode: LaunchMode.externalApplication,
      )) {
        throw Exception('Could not launch browser');
      }

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

      await _connector.disconnectService(serviceName, user['id'], token);

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
    final colors = context.appColors;
    final connectedServices = List<ServiceModel>.from(_visibleConnected);
    _applySort(connectedServices);

    return Scaffold(
      backgroundColor: theme.colorScheme.background,
      body: SafeArea(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
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
                      color: colors.darkGrey,
                    ),
                  ),
                ],
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 24.0),
              child: Column(
                children: [
                  Row(
                    children: [
                      Expanded(
                        child: Semantics(
                          label: 'Rechercher un service connecté',
                          textField: true,
                          child: TextField(
                            controller: _searchController,
                            onChanged: _filterServices,
                            decoration: InputDecoration(
                              hintText: 'Rechercher un service connecté...',
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
                      ),
                      const SizedBox(width: 8),
                      IconButton(
                        tooltip: 'Filtres',
                        onPressed: () {
                          setState(() => _showFilters = !_showFilters);
                        },
                        icon: Icon(
                          Icons.filter_list_rounded,
                          color: _showFilters ? colors.deepBlue : colors.darkGrey,
                        ),
                      ),
                    ],
                  ),
                  if (_showFilters) ...[
                    const SizedBox(height: 10),
                    Wrap(
                      spacing: 10,
                      children: [
                        ChoiceChip(
                          label: const Text('A-Z'),
                          selected: _sort == SortOption.alphabetical,
                          onSelected: (_) {
                            setState(() {
                              _sort = SortOption.alphabetical;
                              _filterServices(_searchController.text);
                            });
                          },
                        ),
                        ChoiceChip(
                          label: const Text('Z-A'),
                          selected: _sort == SortOption.reverseAlphabetical,
                          onSelected: (_) {
                            setState(() {
                              _sort = SortOption.reverseAlphabetical;
                              _filterServices(_searchController.text);
                            });
                          },
                        ),
                      ],
                    ),
                  ],
                ],
              ),
            ),
            const SizedBox(height: 24),
            Expanded(
              child: _isLoading
                  ? const Center(child: CircularProgressIndicator())
                  : _error != null
                      ? Center(
                          child: Column(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Icon(
                                Icons.error_outline,
                                size: 64,
                                color: colors.darkGrey,
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
                            physics: const AlwaysScrollableScrollPhysics(),
                            padding:
                                const EdgeInsets.symmetric(horizontal: 24),
                            children: [
                              Row(
                                mainAxisAlignment:
                                    MainAxisAlignment.spaceBetween,
                                children: [
                                  Text(
                                    'Services connectés',
                                    style: theme.textTheme.titleMedium,
                                  ),
                                  TextButton.icon(
                                    onPressed: _availableServices.isEmpty
                                        ? null
                                        : () => _openConnectSheet(context),
                                    icon: const Icon(Icons.add_link),
                                    label: const Text('Connecter'),
                                  ),
                                ],
                              ),
                              const SizedBox(height: 12),
                              if (connectedServices.isNotEmpty)
                                GridView.builder(
                                  shrinkWrap: true,
                                  physics:
                                      const NeverScrollableScrollPhysics(),
                                  gridDelegate:
                                      const SliverGridDelegateWithFixedCrossAxisCount(
                                    crossAxisCount: 2,
                                    childAspectRatio: 1.2,
                                    crossAxisSpacing: 12,
                                    mainAxisSpacing: 12,
                                  ),
                                  itemCount: connectedServices.length,
                                  itemBuilder: (context, index) {
                                    final service = connectedServices[index];
                                    return _ServiceCard(
                                      service: service,
                                  isConnecting:
                                      _connectingService == service.name,
                                  onConnect: () =>
                                      _connectService(service.name),
                                  onDisconnect: () =>
                                      _disconnectService(service.name),
                                  onDetails: () =>
                                      _showServiceDetails(service),
                                  showDisconnect: true,
                                  showConnectPill: false,
                                  showOverflowDisconnect: true,
                                );
                              },
                            )
                              else
                                _EmptyConnectedState(
                                  onConnect: _availableServices.isEmpty
                                      ? null
                                      : () => _openConnectSheet(context),
                                ),
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

  void _applySort(List<ServiceModel> list) {
    if (_sort == SortOption.alphabetical) {
      list.sort((a, b) => a.displayName.compareTo(b.displayName));
    } else {
      list.sort((a, b) => b.displayName.compareTo(a.displayName));
    }
  }

  void _openConnectSheet(BuildContext context) {
    final colors = context.appColors;
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) {
        return DraggableScrollableSheet(
          initialChildSize: 0.8,
          maxChildSize: 0.95,
          minChildSize: 0.6,
          builder: (context, controller) {
            final filteredAvailable = _availableServices
                .where((s) => s.displayName
                    .toLowerCase()
                    .contains(_connectSearchController.text.toLowerCase()))
                .toList();
            _applySort(filteredAvailable);

            return Container(
              decoration: BoxDecoration(
                color: Theme.of(context).scaffoldBackgroundColor,
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
              child: Padding(
                padding: EdgeInsets.only(
                  left: 20,
                  right: 20,
                  top: 12,
                  bottom: MediaQuery.of(context).viewInsets.bottom + 20,
                ),
                child: Column(
                  children: [
                    Container(
                      width: 40,
                      height: 4,
                      decoration: BoxDecoration(
                        color: colors.grey,
                        borderRadius: BorderRadius.circular(10),
                      ),
                    ),
                    const SizedBox(height: 14),
                    Text(
                      'Connecter un service',
                      style: Theme.of(context).textTheme.titleLarge,
                    ),
                    const SizedBox(height: 12),
                    Semantics(
                      label: 'Rechercher un service à connecter',
                      textField: true,
                      child: TextField(
                        controller: _connectSearchController,
                        onChanged: (_) => setState(() {}),
                        decoration: const InputDecoration(
                          prefixIcon: Icon(Icons.search),
                          hintText: 'Rechercher un service à connecter',
                        ),
                      ),
                    ),
                    const SizedBox(height: 14),
                    Expanded(
                      child: filteredAvailable.isEmpty
                          ? Center(
                              child: Text(
                                'Aucun service disponible',
                                style: Theme.of(context).textTheme.bodyMedium,
                              ),
                            )
                          : GridView.builder(
                              controller: controller,
                              gridDelegate:
                                  const SliverGridDelegateWithFixedCrossAxisCount(
                                crossAxisCount: 2,
                                childAspectRatio: 1.2,
                                crossAxisSpacing: 12,
                                mainAxisSpacing: 12,
                              ),
                              itemCount: filteredAvailable.length,
                              itemBuilder: (context, index) {
                                final service = filteredAvailable[index];
                                return _ServiceCard(
                                  service: service,
                                  isConnecting:
                                      _connectingService == service.name,
                                  onConnect: () =>
                                      _connectService(service.name),
                                  onDisconnect: () {},
                                  onDetails: () =>
                                      _showServiceDetails(service),
                                  showDisconnect: false,
                                  showConnectPill: false,
                                  connectOnTap: true,
                                  showEye: true,
                                  showOverflowDisconnect: false,
                                );
                              },
                            ),
                    ),
                  ],
                ),
              ),
            );
          },
        );
      },
    );
  }

  void _showServiceDetails(ServiceModel service) {
    final colors = context.appColors;
    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.transparent,
      builder: (context) {
        return Container(
          decoration: BoxDecoration(
            color: Theme.of(context).scaffoldBackgroundColor,
            borderRadius:
                const BorderRadius.vertical(top: Radius.circular(18)),
            boxShadow: [
              BoxShadow(
                color: colors.grey.withOpacity(0.25),
                blurRadius: 18,
                offset: const Offset(0, -4),
              ),
            ],
          ),
          child: Padding(
            padding: const EdgeInsets.fromLTRB(20, 12, 20, 24),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Center(
                  child: Container(
                    width: 42,
                    height: 4,
                    decoration: BoxDecoration(
                      color: colors.grey,
                      borderRadius: BorderRadius.circular(10),
                    ),
                  ),
                ),
                const SizedBox(height: 12),
                Row(
                  children: [
                    Container(
                      width: 48,
                      height: 48,
                      decoration: BoxDecoration(
                        gradient: LinearGradient(
                          colors: service.gradient,
                          begin: Alignment.topLeft,
                          end: Alignment.bottomRight,
                        ),
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: Center(
                        child: Text(
                          service.badge,
                          style: const TextStyle(
                            color: Colors.white,
                            fontWeight: FontWeight.w700,
                          ),
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          service.displayName,
                          style: Theme.of(context).textTheme.titleMedium,
                        ),
                        Text(
                          service.isConnected ? 'Connecté' : 'Non connecté',
                          style: Theme.of(context).textTheme.bodySmall?.copyWith(
                                color: service.isConnected
                                    ? Colors.green.shade700
                                    : colors.darkGrey,
                              ),
                        ),
                      ],
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                Text(
                  'Actions',
                  style: Theme.of(context).textTheme.titleSmall,
                ),
                const SizedBox(height: 8),
                if (service.actions.isEmpty)
                  Text(
                    'Aucune action disponible.',
                    style: Theme.of(context)
                        .textTheme
                        .bodySmall
                        ?.copyWith(color: colors.darkGrey),
                  )
                else
                  Wrap(
                    spacing: 8,
                    runSpacing: 8,
                    children: service.actions
                        .map(
                          (a) => Chip(
                            label: Text(a),
                            backgroundColor: colors.lightGrey,
                          ),
                        )
                        .toList(),
                  ),
                const SizedBox(height: 16),
                Text(
                  'Réactions',
                  style: Theme.of(context).textTheme.titleSmall,
                ),
                const SizedBox(height: 8),
                if (service.reactions.isEmpty)
                  Text(
                    'Aucune réaction disponible.',
                    style: Theme.of(context)
                        .textTheme
                        .bodySmall
                        ?.copyWith(color: colors.darkGrey),
                  )
                else
                  Wrap(
                    spacing: 8,
                    runSpacing: 8,
                    children: service.reactions
                        .map(
                          (r) => Chip(
                            label: Text(r),
                            backgroundColor: colors.lightGrey,
                          ),
                        )
                        .toList(),
                  ),
              ],
            ),
          ),
        );
      },
    );
  }
}

class _ServiceCard extends StatelessWidget {
  final ServiceModel service;
  final bool isConnecting;
  final VoidCallback onConnect;
  final VoidCallback onDisconnect;
  final bool showDisconnect;
  final VoidCallback? onDetails;
  final bool showConnectPill;
  final bool connectOnTap;
  final bool showEye;
  final bool showOverflowDisconnect;

  const _ServiceCard({
    required this.service,
    required this.isConnecting,
    required this.onConnect,
    required this.onDisconnect,
    this.showDisconnect = true,
    this.onDetails,
    this.showConnectPill = true,
    this.connectOnTap = false,
    this.showEye = true,
    this.showOverflowDisconnect = false,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colors = context.appColors;

    return InkWell(
      borderRadius: BorderRadius.circular(16),
      onTap: connectOnTap ? onConnect : (onDetails ?? onConnect),
      child: Container(
        padding: const EdgeInsets.all(14.0),
        decoration: BoxDecoration(
          gradient: LinearGradient(
            colors: service.gradient,
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
          ),
          borderRadius: BorderRadius.circular(16),
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
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Row(
              children: [
                CircleAvatar(
                  radius: 18,
                  backgroundColor: Colors.white.withOpacity(0.18),
                  child: Text(
                    service.badge,
                    style: const TextStyle(
                      color: Colors.white,
                      fontWeight: FontWeight.w700,
                      fontSize: 14,
                    ),
                  ),
                ),
                const SizedBox(width: 10),
                _statusDot(service.isConnected),
                const Spacer(),
                if (isConnecting)
                  const SizedBox(
                    width: 20,
                    height: 20,
                    child: CircularProgressIndicator(
                      strokeWidth: 2,
                      color: Colors.white,
                    ),
                  )
                else ...[
                  if (service.isConnected &&
                      showDisconnect &&
                      showOverflowDisconnect)
                    _overflowMenu(onDisconnect)
                  else if (service.isConnected &&
                      showDisconnect &&
                      showConnectPill)
                    _pillButton(label: 'Déconnecter', onTap: onDisconnect)
                  else if (showConnectPill)
                    _pillButton(label: 'À connecter', onTap: onConnect),
                  if (showEye) ...[
                    const SizedBox(width: 8),
                    _iconPill(
                      icon: Icons.remove_red_eye_outlined,
                      onTap: onDetails ?? onConnect,
                    ),
                  ],
                ],
              ],
            ),
            Text(
              service.displayName,
              style: theme.textTheme.titleMedium?.copyWith(
                color: Colors.white,
                fontSize: 15,
                fontWeight: FontWeight.w600,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _pillButton({
    required String label,
    required VoidCallback onTap,
  }) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(24),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.16),
          borderRadius: BorderRadius.circular(24),
          border: Border.all(color: Colors.white.withOpacity(0.35)),
        ),
        child: Text(
          label.toUpperCase(),
          style: const TextStyle(
            color: Colors.white,
            fontWeight: FontWeight.w700,
            fontSize: 12,
            letterSpacing: 0.2,
          ),
        ),
      ),
    );
  }

  Widget _iconPill({required IconData icon, VoidCallback? onTap}) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(24),
      child: Container(
        padding: const EdgeInsets.all(7),
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.16),
          shape: BoxShape.circle,
          border: Border.all(color: Colors.white.withOpacity(0.35)),
        ),
        child: Icon(icon, color: Colors.white, size: 18),
      ),
    );
  }

  Widget _statusDot(bool isConnected) {
    final color = isConnected ? Colors.greenAccent : Colors.redAccent;
    return Container(
      width: 10,
      height: 10,
      decoration: BoxDecoration(
        color: color,
        shape: BoxShape.circle,
        boxShadow: [
          BoxShadow(
            color: color.withOpacity(0.4),
            blurRadius: 8,
            offset: const Offset(0, 2),
          ),
        ],
      ),
    );
  }

  Widget _overflowMenu(VoidCallback onDisconnect) {
    return PopupMenuButton<String>(
      onSelected: (value) {
        if (value == 'disconnect') {
          onDisconnect();
        }
      },
      color: Colors.black.withOpacity(0.85),
      itemBuilder: (context) => [
        PopupMenuItem(
          value: 'disconnect',
          child: Row(
            children: const [
              Icon(Icons.link_off, size: 18, color: Colors.white),
              SizedBox(width: 8),
              Text(
                'Déconnecter',
                style: TextStyle(color: Colors.white),
              ),
            ],
          ),
        ),
      ],
      icon: const Icon(
        Icons.more_horiz,
        color: Colors.white,
      ),
    );
  }

}

class _EmptyConnectedState extends StatelessWidget {
  final VoidCallback? onConnect;

  const _EmptyConnectedState({this.onConnect});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colors = context.appColors;
    return Column(
      children: [
        const SizedBox(height: 24),
        Icon(Icons.link_off, size: 56, color: colors.darkGrey),
        const SizedBox(height: 12),
        Text(
          'Pas encore de service connecté',
          style: theme.textTheme.titleMedium,
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: 8),
        Text(
          'Connectez un service pour le voir apparaître ici.',
          style:
              theme.textTheme.bodySmall?.copyWith(color: colors.darkGrey),
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: 12),
        if (onConnect != null)
          ElevatedButton.icon(
            onPressed: onConnect,
            icon: const Icon(Icons.add_link),
            label: const Text('Connecter un service'),
          ),
      ],
    );
  }
}
