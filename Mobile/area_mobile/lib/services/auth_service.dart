import 'dart:async';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:url_launcher/url_launcher.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:app_links/app_links.dart';

class AuthService {
  final String baseUrl = dotenv.env['BASE_URL'] ?? '';
  final AppLinks _appLinks = AppLinks();

  Future<String> loginWithGoogle() async {
    final res = await http.get(Uri.parse('$baseUrl/auth/google/url'));
    if (res.statusCode != 200) {
      throw Exception('Failed to get Google auth URL');
    }
    final data = jsonDecode(res.body);
    final authUrl = data['auth_url'];
    if (authUrl == null) throw Exception('Invalid response from backend');

    final completer = Completer<Uri>();
    final sub = _appLinks.uriLinkStream.listen((uri) {
      if (uri.scheme == 'area' && uri.host == 'auth') {
        completer.complete(uri);
      }
    });

    try {
      if (!await launchUrl(Uri.parse(authUrl), mode: LaunchMode.externalApplication)) {
        throw Exception('Cannot open browser for Google login');
      }

      final redirected = await completer.future;
      final code = redirected.queryParameters['code'];
      if (code == null) {
        throw Exception('No authorization code in redirect URI');
      }

      final response = await http.post(
        Uri.parse('$baseUrl/auth/google/callback'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({'code': code}),
      );

      if (response.statusCode != 200) {
        throw Exception('Google login failed (${response.statusCode})');
      }

      final body = jsonDecode(response.body);
      return body['jwt'];
    } finally {
      await sub.cancel();
    }
  }
}