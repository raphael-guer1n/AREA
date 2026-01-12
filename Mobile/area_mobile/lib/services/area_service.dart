import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../models/area_model.dart';

class AreaService {
  final String baseUrl;
  final _storage = const FlutterSecureStorage();

  AreaService({required this.baseUrl});

  Future<Map<String, dynamic>> createEvent(CreateEventRequest reqBody) async {
    final token = await _storage.read(key: 'jwt_token');
    if (token == null || token.isEmpty) {
      throw Exception('Missing JWT token');
    }
    final url = Uri.parse('$baseUrl/area_area_api/createEvent');
    final response = await http.post(
      url,
      headers: {
        'Authorization': 'Bearer $token',
        'Content-Type': 'application/json',
      },
      body: jsonEncode(reqBody.toJson()),
    );

    final data = jsonDecode(response.body);
    if (response.statusCode != 200 && response.statusCode != 202) {
      throw Exception(data['error'] ?? 'Failed to create event');
    }
    return data;
  }
}