import 'package:flutter/material.dart';

class ServiceModel {
  final String name;
  final bool isConnected;
  final List<String> actions;
  final List<String> reactions;

  ServiceModel({
    required this.name,
    required this.isConnected,
    this.actions = const [],
    this.reactions = const [],
  });

  factory ServiceModel.fromJson(Map<String, dynamic> json) {
    final actions = (json['actions'] as List?)?.cast<String>() ?? [];
    final reactions = (json['reactions'] as List?)?.cast<String>() ?? [];
    return ServiceModel(
      name: json['provider'] as String,
      isConnected: json['is_logged'] as bool? ?? false,
      actions: actions,
      reactions: reactions,
    );
  }

  String get displayName {
    return name.isNotEmpty
        ? name[0].toUpperCase() + name.substring(1)
        : name;
  }

  String get badge {
    final parts = displayName.split(RegExp(r'[\s_-]+')).where((p) => p.isNotEmpty).toList();
    if (parts.length >= 2) {
      return (parts[0][0] + parts[1][0]).toUpperCase();
    }
    return displayName.length >= 2
        ? displayName.substring(0, 2).toUpperCase()
        : displayName.toUpperCase();
  }

  List<Color> get gradient {
    const palette = [
      [Color(0xFF002642), Color(0xFF0A1A2F)],
      [Color(0xFF840032), Color(0xFFA11248)],
      [Color(0xFFE59500), Color(0xFFF5AE32)],
      [Color(0xFFE5DADA), Color(0xFFF4EDED)],
      [Color(0xFF02040F), Color(0xFF112A46)],
    ];
    final hash = name.hashCode.abs();
    final index = hash % palette.length;
    return palette[index];
  }
}
