import 'dart:convert';
import 'package:http/http.dart' as http;
import '../services/config_service.dart';

class ServiceCatalogService {
  Future<List<String>> fetchServices() async {
    final baseUrl = await ConfigService.getBaseUrl();
    final url = Uri.parse('$baseUrl/area_service_api/providers/services');
    try {
      final response = await http.get(url);
      final body = jsonDecode(response.body);
      if (response.statusCode == 200 &&
          body['success'] == true &&
          body['data']?['services'] is List) {
        final services = List<String>.from(body['data']['services']);
        return services.where((s) => s.trim().isNotEmpty).toList();
      }
      throw Exception(body['error'] ?? 'Failed to load services');
    } catch (e) {
      throw Exception('Error fetching services: $e');
    }
  }
}