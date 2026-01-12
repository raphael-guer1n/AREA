import 'dart:async';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:url_launcher/url_launcher.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:app_links/app_links.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:flutter/foundation.dart';
import '../services/config_service.dart';

class AuthService {
  final String baseUrl = dotenv.env['BASE_URL'] ?? 'http://localhost:8080';

  final AppLinks _appLinks = AppLinks();
  final _storage = const FlutterSecureStorage();

  static const String redirectUri =
      'https://nonbeatifically-stridulatory-denver.ngrok-free.dev/oauth2/callback';

  Future<void> _saveToken(String token) async {
    await _storage.write(key: 'jwt_token', value: token);
    debugPrint('🔐 Saved token: $token');
  }

  Future<String?> getToken() async => _storage.read(key: 'jwt_token');
  Future<void> clearToken() async => _storage.delete(key: 'jwt_token');

  Future<Map<String, dynamic>> loginWithEmail(
      String email, String password) async {
    try {
      final baseUrl = await ConfigService.getBaseUrl();
      final url = Uri.parse("$baseUrl/area_auth_api/auth/login");
      debugPrint("[AUTH] POST $url");

      final response = await http.post(
        url,
        headers: {"Content-Type": "application/json"},
        body: jsonEncode({
          "emailOrUsername": email,
          "password": password,
        }),
      );

      final body = jsonDecode(response.body);
      if (response.statusCode == 200) {
        final token = body["data"]?["token"] ?? body["token"];
        final user = body["data"]?["user"] ?? body["user"];
        if (token != null) {
          await _saveToken(token);
          return {"token": token, "user": user};
        }
        throw Exception("Invalid response format");
      }
      throw Exception(body["error"] ?? "Invalid credentials");
    } catch (e) {
      throw Exception("Login error: $e");
    }
  }

  Future<Map<String, dynamic>> register({
    required String name,
    required String email,
    required String password,
  }) async {
    try {
      final baseUrl = await ConfigService.getBaseUrl();
      final url = Uri.parse("$baseUrl/area_auth_api/auth/register");
      debugPrint("[AUTH] POST $url");

      final response = await http.post(
        url,
        headers: {"Content-Type": "application/json"},
        body: jsonEncode({
          "username": name,
          "email": email,
          "password": password,
        }),
      );

      final body = jsonDecode(response.body);
      if (response.statusCode == 200 || response.statusCode == 201) {
        final token = body["data"]?["token"] ?? body["token"];
        final user = body["data"]?["user"] ?? body["user"];
        if (token != null) {
          await _saveToken(token);
          return {"token": token, "user": user};
        }
        throw Exception("Invalid response format");
      }
      throw Exception(body["error"] ?? "Registration failed");
    } catch (e) {
      throw Exception("Registration error: $e");
    }
  }

  Future<Map<String, dynamic>> loginWithGoogleWithoutUser() async {
    try {
      debugPrint('🌐 Starting Google login');
      final baseUrl = await ConfigService.getBaseUrl();
      final encodedRedirect = Uri.encodeComponent(redirectUri);
      final fullUrl =
          '$baseUrl/area_auth_api/loginwith?provider=google&callback_url=$encodedRedirect&platform=android';

      final response = await http.get(
        Uri.parse(fullUrl),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode != 200) {
        throw Exception('Backend error: ${response.statusCode}');
      }

      final body = jsonDecode(response.body);
      if (body['success'] != true || body['data'] == null) {
        throw Exception('Backend response invalid: ${response.body}');
      }

      String authUrl = body['data']['auth_url']
          .replaceAll(r'\u0026', '&')
          .replaceAll('\u0026', '&');
      debugPrint('🔗 Launching Google OAuth: $authUrl');

      final completer = Completer<Uri>();
      final sub = _appLinks.uriLinkStream.listen((uri) {
        debugPrint('[DEEP LINK] OAuth callback: $uri');
        if (uri.scheme == 'area' && uri.host == 'auth') {
          completer.complete(uri);
        }
      });

      final ok =
          await launchUrl(Uri.parse(authUrl), mode: LaunchMode.externalApplication);
      if (!ok) throw Exception('Cannot open browser');

      final redirected =
          await completer.future.timeout(const Duration(minutes: 5));
      await sub.cancel();

      final tokenParam = redirected.queryParameters['token'];
      if (tokenParam == null || tokenParam.isEmpty) {
        throw Exception('Missing token in callback');
      }

      await _saveToken(tokenParam);
      return {
        'token': tokenParam,
        'provider': body['data']['provider'] ?? 'google',
      };
    } catch (e) {
      throw Exception('Google login failed: $e');
    }
  }

  Future<Map<String, dynamic>> fetchCurrentUser() async {
    try {
      final baseUrl = await ConfigService.getBaseUrl();
      final token = await _storage.read(key: 'jwt_token');
      if (token == null || token.isEmpty) {
        throw Exception('Missing token');
      }

      final url = Uri.parse('$baseUrl/area_auth_api/auth/me');
      debugPrint('📡 GET $url');
      final response = await http.get(url, headers: {
        'Authorization': 'Bearer $token',
        'Content-Type': 'application/json',
      });

      final body = jsonDecode(response.body);
      if (response.statusCode == 200 &&
          body['success'] == true &&
          body['data']?['user'] != null) {
        return Map<String, dynamic>.from(body['data']['user']);
      }
      throw Exception('Fetch user failed: ${response.body}');
    } catch (e) {
      throw Exception('Fetch user error: $e');
    }
  }

  Future<void> logout() async {
    await clearToken();
  }
}