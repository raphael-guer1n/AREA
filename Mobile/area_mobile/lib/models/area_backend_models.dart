class InputFieldDto {
  final String name;
  final String value;

  InputFieldDto({
    required this.name,
    required this.value,
  });

  factory InputFieldDto.fromJson(Map<String, dynamic> json) {
    return InputFieldDto(
      name: json['name'] as String? ?? '',
      value: json['value'] as String? ?? '',
    );
  }

  Map<String, dynamic> toJson() => {
        'name': name,
        'value': value,
      };
}

class AreaActionDto {
  final bool active;
  final int id;
  final String provider;
  final String service;
  final String title;
  final String type;
  final List<InputFieldDto> input;

  AreaActionDto({
    required this.active,
    required this.id,
    required this.provider,
    required this.service,
    required this.title,
    required this.type,
    required this.input,
  });

  factory AreaActionDto.fromJson(Map<String, dynamic> json) {
    return AreaActionDto(
      active: json['active'] as bool? ?? false,
      id: (json['id'] as num?)?.toInt() ?? 0,
      provider: json['provider'] as String? ?? '',
      service: json['service'] as String? ?? '',
      title: json['title'] as String? ?? '',
      type: json['type'] as String? ?? '',
      input: ((json['input'] as List?) ?? [])
          .map((e) => InputFieldDto.fromJson(Map<String, dynamic>.from(e)))
          .toList(),
    );
  }

  Map<String, dynamic> toJson() => {
        'active': active,
        'id': id,
        'provider': provider,
        'service': service,
        'title': title,
        'type': type,
        'input': input.map((e) => e.toJson()).toList(),
      };
}

class AreaReactionDto {
  final int id;
  final String provider;
  final String service;
  final String title;
  final List<InputFieldDto> input;

  AreaReactionDto({
    required this.id,
    required this.provider,
    required this.service,
    required this.title,
    required this.input,
  });

  factory AreaReactionDto.fromJson(Map<String, dynamic> json) {
    return AreaReactionDto(
      id: (json['id'] as num?)?.toInt() ?? 0,
      provider: json['provider'] as String? ?? '',
      service: json['service'] as String? ?? '',
      title: json['title'] as String? ?? '',
      input: ((json['input'] as List?) ?? [])
          .map((e) => InputFieldDto.fromJson(Map<String, dynamic>.from(e)))
          .toList(),
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'provider': provider,
        'service': service,
        'title': title,
        'input': input.map((e) => e.toJson()).toList(),
      };
}

class AreaDto {
  final int id;
  final String name;
  final bool active;
  final int userId;
  final List<AreaActionDto> actions;
  final List<AreaReactionDto> reactions;

  AreaDto({
    required this.id,
    required this.name,
    required this.active,
    required this.userId,
    required this.actions,
    required this.reactions,
  });

  factory AreaDto.fromJson(Map<String, dynamic> json) {
    return AreaDto(
      id: (json['id'] as num?)?.toInt() ?? 0,
      name: json['name'] as String? ?? '',
      active: json['active'] as bool? ?? false,
      userId: (json['user_id'] as num?)?.toInt() ?? 0,
      actions: ((json['actions'] as List?) ?? [])
          .map((e) => AreaActionDto.fromJson(Map<String, dynamic>.from(e)))
          .toList(),
      reactions: ((json['reactions'] as List?) ?? [])
          .map((e) => AreaReactionDto.fromJson(Map<String, dynamic>.from(e)))
          .toList(),
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'name': name,
        'active': active,
        'user_id': userId,
        'actions': actions.map((e) => e.toJson()).toList(),
        'reactions': reactions.map((e) => e.toJson()).toList(),
      };

  AreaDto copyWith({bool? active}) {
    return AreaDto(
      id: id,
      name: name,
      active: active ?? this.active,
      userId: userId,
      actions: actions,
      reactions: reactions,
    );
  }
}