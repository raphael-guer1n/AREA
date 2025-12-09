import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/service_model.dart';  // Add this import

class ServiceConnector {
  final String baseUrl;

  ServiceConnector({required this.baseUrl});

  Future<List<ServiceModel>> fetchServices(int userId) async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/oauth2/providers/$userId'),
      );

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        if (data['success'] == true && data['data'] != null) {
          final providers = data['data']['providers'] as List;
          return providers
              .map((p) => ServiceModel.fromJson(p))
              .toList();
        }
      }
      throw Exception('Failed to load services');
    } catch (e) {
      throw Exception('Error fetching services: $e');
    }
  }

  Future<String> getAuthUrl(String serviceName, int userId) async {
    try {
      final response = await http.get(
        Uri.parse(
          '$baseUrl/oauth2/authorize?provider=$serviceName&user_id=$userId',
        ),
      );

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        if (data['success'] == true && data['data'] != null) {
          return data['data']['auth_url'];
        }
      }
      throw Exception('Failed to get auth URL');
    } catch (e) {
      throw Exception('Error getting auth URL: $e');
    }
  }

  Future<void> disconnectService(
    String serviceName,
    int userId,
    String token,
  ) async {
    // Implement disconnect logic based on your backend API
    // This is a placeholder
    try {
      // You might need to add a disconnect endpoint to your backend
      throw UnimplementedError('Disconnect endpoint not implemented yet');
    } catch (e) {
      throw Exception('Error disconnecting service: $e');
    }
  }
}