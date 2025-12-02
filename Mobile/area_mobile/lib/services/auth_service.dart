import 'dart:convert';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:http/http.dart' as http;
import 'oauth_helper.dart';

class AuthService {
  final String baseUrl = dotenv.env['BASE_URL'] ?? '';

  Future<String> loginWithGoogle() async {
    final urlResponse = await http.get(Uri.parse('$baseUrl/auth/google/url'));
    if (urlResponse.statusCode != 200) {
      throw Exception('Error retrieving Google login URL');
    }

    final data = jsonDecode(urlResponse.body);
    final authUrl = data['auth_url'];

    final callback = await OAuthHelper.launchAuthFlow(
      url: authUrl,
      callbackScheme: 'area',
    );

    final uri = Uri.parse(callback);
    final code = uri.queryParameters['code'];

    if (code == null) {
      throw Exception('OAuth2 callback missing code');
    }

    final response = await http.post(
      Uri.parse('$baseUrl/auth/google/callback'),
      body: jsonEncode({'code': code}),
      headers: {'Content-Type': 'application/json'},
    );

    if (response.statusCode != 200) {
      throw Exception('Google login failed');
    }

    final body = jsonDecode(response.body);
    return body['jwt'];
  }
}