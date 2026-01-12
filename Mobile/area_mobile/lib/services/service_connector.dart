import "dart:convert";
import "package:flutter_secure_storage/flutter_secure_storage.dart";
import "package:http/http.dart" as http;
import "../models/service_model.dart";
import "../services/config_service.dart";

class ServiceConnector {
  final _storage = const FlutterSecureStorage();

  Future<List<ServiceModel>> fetchServices(int userId) async {
    try {
      final token = await _storage.read(key: "jwt_token");
      if (token == null || token.isEmpty) {
        throw Exception("JWT token missing — please log in again.");
      }

      final baseUrl = await ConfigService.getBaseUrl();
      final url = Uri.parse("$baseUrl/area_auth_api/oauth2/providers/$userId");

      final response = await http.get(
        url,
        headers: {"Authorization": "Bearer $token"},
      );

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        if (data["success"] == true && data["data"] != null) {
          final providers = data["data"]["providers"] as List;
          return providers.map((p) => ServiceModel.fromJson(p)).toList();
        }
      }

      throw Exception("Failed to load services (${response.statusCode})");
    } catch (e) {
      throw Exception("Error fetching services: $e");
    }
  }

  Future<String> getAuthUrl(String serviceName, int userId) async {
    try {
      final token = await _storage.read(key: "jwt_token");
      if (token == null || token.isEmpty) {
        throw Exception("JWT token missing — please log in first.");
      }

      final baseUrl = await ConfigService.getBaseUrl();
      const redirectUri =
          "https://nonbeatifically-stridulatory-denver.ngrok-free.dev/oauth2/callback";
      final encodedRedirect = Uri.encodeComponent(redirectUri);

      final uri = Uri.parse(
        "$baseUrl/area_auth_api/oauth2/authorize"
        "?provider=$serviceName&user_id=$userId"
        "&callback_url=$encodedRedirect&platform=android",
      );

      final response = await http.get(
        uri,
        headers: {
          "Authorization": "Bearer $token",
          "Content-Type": "application/json",
        },
      );

      if (response.statusCode != 200) {
        throw Exception("Failed to get auth URL (${response.statusCode})");
      }

      final data = jsonDecode(response.body);
      if (data["success"] != true || data["data"] == null) {
        throw Exception(data["error"] ?? "Invalid backend response");
      }

      String authUrl = (data["data"]["auth_url"] as String)
          .replaceAll("\\u0026", "&")
          .replaceAll("\u0026", "&");

      return authUrl;
    } catch (e) {
      throw Exception("Error getting auth URL: $e");
    }
  }

  Future<void> disconnectService(
    String serviceName,
    int userId,
    String token,
  ) async {
    throw UnimplementedError("Disconnect endpoint not implemented yet.");
  }
}