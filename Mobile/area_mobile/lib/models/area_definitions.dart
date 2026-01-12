import 'package:flutter/material.dart';

enum AreaFieldType {
  text,
  number,
  date,
}

class AreaFieldDefinition {
  final String name;
  final AreaFieldType type;
  final String label;
  final bool required;
  final String? defaultValue;
  final bool isPrivate;

  const AreaFieldDefinition({
    required this.name,
    required this.type,
    required this.label,
    this.required = false,
    this.defaultValue,
    this.isPrivate = false,
  });
}

class AreaActionDefinition {
  final String id;
  final String title;
  final String label;
  final String type;
  final List<AreaFieldDefinition> fields;

  const AreaActionDefinition({
    required this.id,
    required this.title,
    required this.label,
    required this.type,
    this.fields = const [],
  });
}

class AreaReactionDefinition {
  final String id;
  final String title;
  final String label;
  final List<AreaFieldDefinition> fields;

  const AreaReactionDefinition({
    required this.id,
    required this.title,
    required this.label,
    this.fields = const [],
  });
}

class AreaServiceDefinition {
  final String id;
  final String name;
  final List<AreaActionDefinition> actions;
  final List<AreaReactionDefinition> reactions;
  final bool connected;

  const AreaServiceDefinition({
    required this.id,
    required this.name,
    this.actions = const [],
    this.reactions = const [],
    this.connected = false,
  });

  AreaServiceDefinition copyWith({
    bool? connected,
  }) {
    return AreaServiceDefinition(
      id: id,
      name: name,
      actions: actions,
      reactions: reactions,
      connected: connected ?? this.connected,
    );
  }

  String get displayName {
    if (name.isEmpty) return name;
    return name[0].toUpperCase() + name.substring(1);
  }
}

class AreaGradient {
  final Color from;
  final Color to;

  const AreaGradient({
    required this.from,
    required this.to,
  });
}

class CreatedArea {
  final String id;
  final String name;
  final String summary;
  final String startTime;
  final String endTime;
  final int delay;
  final String actionService;
  final String reactionService;
  final String actionName;
  final String reactionName;
  final AreaGradient gradient;

  const CreatedArea({
    required this.id,
    required this.name,
    required this.summary,
    required this.startTime,
    required this.endTime,
    required this.delay,
    required this.actionService,
    required this.reactionService,
    required this.actionName,
    required this.reactionName,
    required this.gradient,
  });
}
