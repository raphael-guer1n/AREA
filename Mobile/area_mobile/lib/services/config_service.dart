import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';

class ConfigService {
  static const _keyServerIp = 'server_ip';
  static final FlutterSecureStorage _storage = const FlutterSecureStorage();

  /// Get currently stored server IP or fall back to .env BASE_URL
  static Future<String> getBaseUrl() async {
    final customIp = await _storage.read(key: _keyServerIp);
    if (customIp != null && customIp.trim().isNotEmpty) {
      return 'http://$customIp:8080';
    }
    return dotenv.env['BASE_URL'] ?? 'http://localhost:8080';
  }

  /// Save new server IP (just the raw IP, not with http://)
  static Future<void> setServerIp(String ip) async {
    await _storage.write(key: _keyServerIp, value: ip.trim());
  }

  /// Read the saved raw IP value (not the full URL)
  static Future<String?> getServerIp() async {
    return await _storage.read(key: _keyServerIp);
  }

  /// Clear the IP override (revert to .env BASE_URL)
  static Future<void> clearServerIp() async {
    await _storage.delete(key: _keyServerIp);
  }
}