import 'dart:convert';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:http/http.dart' as http;
import '../models/area_backend_models.dart';
import '../services/config_service.dart';

class SaveAreaResult {
  final bool success;
  final String? message;
  final List<String> missingProviders;

  SaveAreaResult({
    required this.success,
    this.message,
    this.missingProviders = const [],
  });
}

class AreaService {
  final _storage = const FlutterSecureStorage();

  Future<String> _requireToken() async {
    final token = await _storage.read(key: 'jwt_token');
    if (token == null || token.isEmpty) {
      throw Exception('Missing JWT token');
    }
    return token;
  }

  Future<List<AreaDto>> getAreas() async {
    final token = await _requireToken();
    final baseUrl = await ConfigService.getBaseUrl();
    final url = Uri.parse('$baseUrl/area_area_api/getAreas');

    final response = await http.get(url, headers: {
      'Authorization': 'Bearer $token',
      'Content-Type': 'application/json',
    });

    final body = jsonDecode(response.body) as Map<String, dynamic>;
    if (response.statusCode == 200 && body['success'] == true) {
      final list = (body['data'] as List?) ?? [];
      return list
          .map((e) => AreaDto.fromJson(Map<String, dynamic>.from(e)))
          .toList();
    }
    throw Exception(body['error'] ?? 'Failed to load areas');
  }

  Future<SaveAreaResult> saveArea(AreaDto area) async {
    final token = await _requireToken();
    final baseUrl = await ConfigService.getBaseUrl();
    final url = Uri.parse('$baseUrl/area_area_api/saveArea');

    final response = await http.post(
      url,
      headers: {
        'Authorization': 'Bearer $token',
        'Content-Type': 'application/json',
      },
      body: jsonEncode(area.toJson()),
    );

    final body = jsonDecode(response.body) as Map<String, dynamic>;

    if (response.statusCode == 200 && body['success'] == true) {
      final missing = (body['missing_providers'] as List?)
              ?.map((e) => e.toString())
              .toList() ??
          const <String>[];

      return SaveAreaResult(
        success: true,
        message: body['message'] as String?,
        missingProviders: missing,
      );
    }

    throw Exception(body['error'] ?? 'Failed to save area');
  }

  Future<void> activateArea(int areaId) async {
    final token = await _requireToken();
    final baseUrl = await ConfigService.getBaseUrl();
    final url = Uri.parse('$baseUrl/area_area_api/activateArea');

    final response = await http.post(
      url,
      headers: {
        'Authorization': 'Bearer $token',
        'Content-Type': 'application/json',
      },
      body: jsonEncode({'area_id': areaId}),
    );

    final body = jsonDecode(response.body) as Map<String, dynamic>;
    if (response.statusCode == 200) {
      if (body['success'] == false) {
        throw Exception(body['error'] ?? body['message'] ?? 'Activate failed');
      }
      return;
    }
    throw Exception(body['error'] ?? 'Activate failed');
  }

  Future<void> deactivateArea(int areaId) async {
    final token = await _requireToken();
    final baseUrl = await ConfigService.getBaseUrl();
    final url = Uri.parse('$baseUrl/area_area_api/deactivateArea');

    final response = await http.post(
      url,
      headers: {
        'Authorization': 'Bearer $token',
        'Content-Type': 'application/json',
      },
      body: jsonEncode({'area_id': areaId}),
    );

    final body = jsonDecode(response.body) as Map<String, dynamic>;
    if (response.statusCode == 200) {
      if (body['success'] == false) {
        throw Exception(
          body['error'] ?? body['message'] ?? 'Deactivate failed',
        );
      }
      return;
    }
    throw Exception(body['error'] ?? 'Deactivate failed');
  }
}