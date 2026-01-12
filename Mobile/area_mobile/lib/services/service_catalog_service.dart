import 'dart:convert';
import 'package:http/http.dart' as http;

class ServiceCatalogService {
  final String baseUrl;

  ServiceCatalogService({required this.baseUrl});

  Future<List<String>> fetchServices() async {
    final url = Uri.parse('$baseUrl/area_service_api/providers/services');
    try {
      final response = await http.get(url);
      final body = jsonDecode(response.body);
      if (response.statusCode == 200 &&
          body['success'] == true &&
          body['data']?['services'] is List) {
        final services = List<String>.from(body['data']['services']);
        return services.where((service) => service.trim().isNotEmpty).toList();
      }
      throw Exception(body['error'] ?? 'Failed to load services');
    } catch (e) {
      throw Exception('Error fetching services: $e');
    }
  }
}
