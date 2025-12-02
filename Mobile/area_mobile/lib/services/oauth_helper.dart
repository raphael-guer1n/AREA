import 'package:flutter_web_auth_2/flutter_web_auth_2.dart';

class OAuthHelper {
  static Future<String> launchAuthFlow({
    required String url,
    required String callbackScheme,
  }) async {
    try {
      final result = await FlutterWebAuth2.authenticate(
        url: url,
        callbackUrlScheme: callbackScheme,
      );
      return result;
    } catch (e) {
      throw Exception('OAuth2 flow failed: $e');
    }
  }
}