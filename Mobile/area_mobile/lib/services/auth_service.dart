import 'dart:async';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:url_launcher/url_launcher.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:app_links/app_links.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class AuthService {
  final String baseUrl = dotenv.env['BASE_URL'] ?? 'http://10.0.2.2:8083';
  final AppLinks _appLinks = AppLinks();
  final _storage = const FlutterSecureStorage();

  Future<void> _saveToken(String token) async {
    await _storage.write(key: 'jwt_token', value: token);
  }

  Future<String?> getToken() async {
    return await _storage.read(key: 'jwt_token');
  }

  Future<void> clearToken() async {
    await _storage.delete(key: 'jwt_token');
  }

  Future<Map<String, dynamic>> loginWithEmail(
      String email, String password) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/auth/login'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'emailOrUsername': email,
          'password': password,
        }),
      );

      if (response.statusCode == 200) {
        final body = jsonDecode(response.body);

        if (body['success'] == true && body['data'] != null) {
          final token = body['data']['token'];
          final user = body['data']['user'];
          await _saveToken(token);
          return {'token': token, 'user': user};
        } else if (body['token'] != null || body['jwt'] != null) {
          final token = body['token'] ?? body['jwt'];
          await _saveToken(token);
          return {'token': token, 'user': body['user']};
        }

        throw Exception('Format de réponse invalide');
      } else {
        final body = jsonDecode(response.body);
        throw Exception(body['error'] ?? 'Identifiants incorrects');
      }
    } catch (e) {
      throw Exception('Erreur de connexion: $e');
    }
  }

  Future<Map<String, dynamic>> register({
    required String name,
    required String email,
    required String password,
  }) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/auth/register'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'username': name,
          'email': email,
          'password': password,
        }),
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        final body = jsonDecode(response.body);

        if (body['success'] == true && body['data'] != null) {
          final token = body['data']['token'];
          final user = body['data']['user'];
          await _saveToken(token);
          return {'token': token, 'user': user};
        } else if (body['token'] != null || body['jwt'] != null) {
          final token = body['token'] ?? body['jwt'];
          await _saveToken(token);
          return {'token': token, 'user': body['user']};
        }

        throw Exception('Format de réponse invalide');
      } else {
        final body = jsonDecode(response.body);
        throw Exception(
            body['error'] ?? 'Erreur lors de la création du compte');
      }
    } catch (e) {
      throw Exception('Erreur d\'inscription: $e');
    }
  }

  Future<Map<String, dynamic>> _handleOAuthFlow(String provider) async {
    final res = await http.get(
      Uri.parse('$baseUrl/auth/oauth2/authorize?provider=$provider'),
    );

    if (res.statusCode != 200) {
      throw Exception('Impossible d\'obtenir l\'URL d\'authentification');
    }

    final data = jsonDecode(res.body);

    String? authUrl;
    if (data['success'] == true && data['data'] != null) {
      authUrl = data['data']['auth_url'];
    } else {
      authUrl = data['auth_url'];
    }

    if (authUrl == null) {
      throw Exception('URL d\'authentification invalide');
    }

    final completer = Completer<Uri>();
    final sub = _appLinks.uriLinkStream.listen((uri) {
      if (uri.scheme == 'area' && uri.host == 'auth') {
        completer.complete(uri);
      }
    });

    try {
      if (!await launchUrl(
        Uri.parse(authUrl),
        mode: LaunchMode.externalApplication,
      )) {
        throw Exception('Impossible d\'ouvrir le navigateur');
      }

      final redirected =
          await completer.future.timeout(const Duration(minutes: 5));
      final code = redirected.queryParameters['code'];
      final state = redirected.queryParameters['state'];

      if (code == null) {
        throw Exception('Code d\'autorisation manquant');
      }

      final response = await http.get(
        Uri.parse(
            '$baseUrl/auth/oauth2/callback?code=$code&state=${state ?? ""}'),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode != 200) {
        throw Exception('Échec de l\'authentification');
      }

      final body = jsonDecode(response.body);

      String? token;
      Map<String, dynamic>? user;

      if (body['success'] == true && body['data'] != null) {
        token = body['data']['access_token'];
        user = body['data']['user_info'];
      } else {
        token = body['access_token'] ?? body['token'] ?? body['jwt'];
        user = body['user_info'] ?? body['user'];
      }

      if (token == null) {
        throw Exception('Token non trouvé');
      }

      await _saveToken(token);
      return {'token': token, 'user': user};
    } finally {
      await sub.cancel();
    }
  }

  Future<Map<String, dynamic>> loginWithGoogle() =>
      _handleOAuthFlow('google');
  Future<Map<String, dynamic>> loginWithApple() => _handleOAuthFlow('apple');
  Future<Map<String, dynamic>> loginWithFacebook() =>
      _handleOAuthFlow('facebook');

  Future<void> logout() async {
    await clearToken();
  }
}