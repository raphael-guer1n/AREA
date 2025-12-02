import 'dart:convert';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:http/http.dart' as http;
import 'oauth_helper.dart';

class ServiceConnector {
  final String baseUrl = dotenv.env['BASE_URL'] ?? '';

  Future<void> connectToService(String serviceName, String jwtToken) async {
    final response = await http.get(
      Uri.parse('$baseUrl/auth/$serviceName/url'),
      headers: {'Authorization': 'Bearer $jwtToken'},
    );

    if (response.statusCode != 200) {
      throw Exception('Cannot get OAuth2 URL for $serviceName');
    }

    final data = jsonDecode(response.body);
    final authUrl = data['auth_url'];

    await OAuthHelper.launchAuthFlow(
      url: authUrl,
      callbackScheme: 'area',
    );

  }
}