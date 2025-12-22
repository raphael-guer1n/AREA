import 'dart:async';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:url_launcher/url_launcher.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:app_links/app_links.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class AuthService {
  final String baseUrl = dotenv.env['BASE_URL'] ??
      'https://nonbeatifically-stridulatory-denver.ngrok-free.dev';
  final AppLinks _appLinks = AppLinks();
  final _storage = const FlutterSecureStorage();

  static const String redirectUri =
      'https://nonbeatifically-stridulatory-denver.ngrok-free.dev/oauth2/callback';

  Future<void> _saveToken(String token) async {
    await _storage.write(key: 'jwt_token', value: token);
  }

  Future<String?> getToken() async => _storage.read(key: 'jwt_token');
  Future<void> clearToken() async => _storage.delete(key: 'jwt_token');

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
        final token = body['data']?['token'] ?? body['token'];
        final user = body['data']?['user'] ?? body['user'];

        if (token != null) {
          await _saveToken(token);
          return {'token': token, 'user': user};
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
        final token = body['data']?['token'] ?? body['token'];
        final user = body['data']?['user'] ?? body['user'];

        if (token != null) {
          await _saveToken(token);
          return {'token': token, 'user': user};
        }
        throw Exception('Format de réponse invalide');
      } else {
        final body = jsonDecode(response.body);
        throw Exception(body['error'] ?? 'Erreur lors de la création du compte');
      }
    } catch (e) {
      throw Exception('Erreur d\'inscription: $e');
    }
  }

  Future<Map<String, dynamic>> _handleOAuthFlow(
    String provider, {
    required int userId,
  }) async {
    final token = await _storage.read(key: 'jwt_token');

    if (token == null || token.isEmpty) {
      throw Exception('No JWT token found in secure storage. Login first.');
    }
    if (userId <= 0) {
      throw Exception('Invalid user ID ($userId)');
    }

    final encodedRedirect = Uri.encodeComponent(redirectUri);
    final fullUrl =
        '$baseUrl/oauth2/authorize?provider=$provider&user_id=$userId'
        '&callback_url=$encodedRedirect&platform=android';

    final completer = Completer<Uri>();
    final sub = _appLinks.uriLinkStream.listen((uri) {
      if (uri.scheme == 'area' && uri.host == 'auth') {
        completer.complete(uri);
      }
    });

    try {
      final ok = await launchUrl(Uri.parse(fullUrl),
          mode: LaunchMode.externalApplication);
      if (!ok) throw Exception('Impossible d\'ouvrir le navigateur');

      final redirected =
          await completer.future.timeout(const Duration(minutes: 5));

      final tokenParam = redirected.queryParameters['token'];
      final providerParam = redirected.queryParameters['provider'];

      if (tokenParam == null || tokenParam.isEmpty) {
        throw Exception('Jeton manquant dans la redirection');
      }

      await _saveToken(tokenParam);

      return {
        'token': tokenParam,
        'provider': providerParam ?? provider,
      };
    } finally {
      await sub.cancel();
    }
  }

  Future<Map<String, dynamic>> loginWithGoogle({required int userId}) =>
      _handleOAuthFlow('google', userId: userId);

  Future<Map<String, dynamic>> loginWithApple({required int userId}) =>
      _handleOAuthFlow('apple', userId: userId);

  Future<Map<String, dynamic>> loginWithFacebook({required int userId}) =>
      _handleOAuthFlow('facebook', userId: userId);

  Future<void> logout() async => clearToken();
}