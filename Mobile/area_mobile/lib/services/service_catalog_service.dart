import 'dart:convert';
import 'package:http/http.dart' as http;

import '../services/config_service.dart';

class ServiceFieldConfigDto {
  final String name;
  final String type;
  final String label;
  final bool required;
  final String defaultValue;

  ServiceFieldConfigDto({
    required this.name,
    required this.type,
    required this.label,
    required this.required,
    required this.defaultValue,
  });

  factory ServiceFieldConfigDto.fromJson(Map<String, dynamic> json) {
    return ServiceFieldConfigDto(
      name: json['name'] as String? ?? '',
      type: json['type'] as String? ?? 'text',
      label: json['label'] as String? ?? '',
      required: json['required'] as bool? ?? false,
      defaultValue: json['default']?.toString() ?? '',
    );
  }
}

class OutputFieldDto {
  final String name;
  final String type;
  final String label;

  OutputFieldDto({
    required this.name,
    required this.type,
    required this.label,
  });

  factory OutputFieldDto.fromJson(Map<String, dynamic> json) {
    return OutputFieldDto(
      name: json['name'] as String? ?? '',
      type: json['type'] as String? ?? 'string',
      label: json['label'] as String? ?? '',
    );
  }
}

class ActionConfigDto {
  final String title;
  final String label;
  final String type;
  final List<ServiceFieldConfigDto> fields;
  final List<OutputFieldDto> outputFields;

  ActionConfigDto({
    required this.title,
    required this.label,
    required this.type,
    required this.fields,
    required this.outputFields,
  });

  factory ActionConfigDto.fromJson(Map<String, dynamic> json) {
    return ActionConfigDto(
      title: json['title'] as String? ?? '',
      label: json['label'] as String? ?? '',
      type: json['type'] as String? ?? '',
      fields: ((json['fields'] as List?) ?? [])
          .map(
            (e) => ServiceFieldConfigDto.fromJson(
              Map<String, dynamic>.from(e),
            ),
          )
          .toList(),
      outputFields: ((json['output_fields'] as List?) ?? [])
          .map((e) => OutputFieldDto.fromJson(Map<String, dynamic>.from(e)))
          .toList(),
    );
  }
}

class ReactionConfigDto {
  final String title;
  final String label;
  final List<ServiceFieldConfigDto> fields;

  ReactionConfigDto({
    required this.title,
    required this.label,
    required this.fields,
  });

  factory ReactionConfigDto.fromJson(Map<String, dynamic> json) {
    return ReactionConfigDto(
      title: json['title'] as String? ?? '',
      label: json['label'] as String? ?? '',
      fields: ((json['fields'] as List?) ?? [])
          .map(
            (e) => ServiceFieldConfigDto.fromJson(
              Map<String, dynamic>.from(e),
            ),
          )
          .toList(),
    );
  }
}

class ServiceConfigDto {
  final String provider;
  final String name;
  final String label;
  final List<ActionConfigDto> actions;
  final List<ReactionConfigDto> reactions;

  ServiceConfigDto({
    required this.provider,
    required this.name,
    required this.label,
    required this.actions,
    required this.reactions,
  });

  factory ServiceConfigDto.fromJson(Map<String, dynamic> json) {
    return ServiceConfigDto(
      provider: json['provider'] as String? ?? '',
      name: json['name'] as String? ?? '',
      label: json['label'] as String? ?? '',
      actions: ((json['actions'] as List?) ?? [])
          .map((e) => ActionConfigDto.fromJson(Map<String, dynamic>.from(e)))
          .toList(),
      reactions: ((json['reactions'] as List?) ?? [])
          .map((e) => ReactionConfigDto.fromJson(Map<String, dynamic>.from(e)))
          .toList(),
    );
  }
}

class ServiceCatalogService {
  Future<List<String>> fetchServiceNames() async {
    final baseUrl = await ConfigService.getBaseUrl();
    final url = Uri.parse('$baseUrl/area_service_api/services/services');

    final response = await http.get(url);
    final body = jsonDecode(response.body) as Map<String, dynamic>;

    if (response.statusCode == 200 &&
        body['success'] == true &&
        body['data']?['services'] is List) {
      final services = List<String>.from(body['data']['services']);
      return services.where((s) => s.trim().isNotEmpty).toList();
    }

    throw Exception(body['error'] ?? 'Failed to load services');
  }

  Future<ServiceConfigDto> fetchServiceConfig(String serviceName) async {
    final baseUrl = await ConfigService.getBaseUrl();
    final url = Uri.parse(
      '$baseUrl/area_service_api/services/service-config?service=$serviceName',
    );

    final response = await http.get(url);
    final body = jsonDecode(response.body) as Map<String, dynamic>;

    if (response.statusCode == 200 && body['success'] == true) {
      final data = Map<String, dynamic>.from(body['data'] as Map);
      return ServiceConfigDto.fromJson(data);
    }

    throw Exception(body['error'] ?? 'Failed to load service config');
  }
}