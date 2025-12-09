class ServiceModel {
  final String name;
  final bool isConnected;

  ServiceModel({
    required this.name,
    required this.isConnected,
  });

  factory ServiceModel.fromJson(Map<String, dynamic> json) {
    return ServiceModel(
      name: json['provider'] as String,
      isConnected: json['is_logged'] as bool? ?? false,
    );
  }

  String get displayName {
    return name[0].toUpperCase() + name.substring(1);
  }

  String get iconPath {
    final iconMap = {
      'google': 'assets/icons/google.png',
      'github': 'assets/icons/github.png',
      'discord': 'assets/icons/discord.png',
      'spotify': 'assets/icons/spotify.png',
      'notion': 'assets/icons/notion.png',
      'slack': 'assets/icons/slack.png',
    };
    return iconMap[name.toLowerCase()] ?? 'assets/icons/default.png';
  }
}