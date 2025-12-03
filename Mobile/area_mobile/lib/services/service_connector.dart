import 'dart:async';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:url_launcher/url_launcher.dart';
import 'package:app_links/app_links.dart';

class ServiceConnector {
  final String baseUrl;
  final AppLinks _appLinks = AppLinks();

  ServiceConnector({required this.baseUrl});

  Future<void> connectToService(String serviceName, String jwtToken) async {
    final res = await http.get(
      Uri.parse('$baseUrl/auth/$serviceName/url'),
      headers: {'Authorization': 'Bearer $jwtToken'},
    );
    if (res.statusCode != 200) {
      throw Exception('Failed to get $serviceName auth URL');
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
      final uri = Uri.parse(authUrl);
      if (!await launchUrl(uri, mode: LaunchMode.externalApplication)) {
        throw Exception('Could not launch browser for $serviceName');
      }

      final redirected = await completer.future;
      final code = redirected.queryParameters['code'];
      if (code == null) throw Exception('Missing authorization code');

      final response = await http.post(
        Uri.parse('$baseUrl/auth/$serviceName/callback'),
        headers: {
          'Authorization': 'Bearer $jwtToken',
          'Content-Type': 'application/json',
        },
        body: jsonEncode({'code': code}),
      );

      if (response.statusCode != 200) {
        throw Exception('Token exchange failed (${response.statusCode})');
      }
    } finally {
      await sub.cancel();
    }
  }
}