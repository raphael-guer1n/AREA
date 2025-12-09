import 'package:flutter/foundation.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../services/auth_service.dart';

class AuthProvider extends ChangeNotifier {
  final AuthService _authService = AuthService();
  final _storage = const FlutterSecureStorage();

  bool _isAuthenticated = false;
  bool _isLoading = true;
  String? _token;
  Map<String, dynamic>? _user;
  String? _error;

  bool get isAuthenticated => _isAuthenticated;
  bool get isLoading => _isLoading;
  String? get token => _token;
  Map<String, dynamic>? get user => _user;
  String? get error => _error;

  AuthProvider() {
    _checkAuthStatus();
  }

  Future<void> _checkAuthStatus() async {
    _isLoading = true;
    notifyListeners();

    try {
      final savedToken = await _authService.getToken();
      if (savedToken != null) {
        _token = savedToken;
        _isAuthenticated = true;
      }
    } catch (e) {
      debugPrint('Error checking auth status: $e');
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<bool> loginWithEmail(String email, String password) async {
    _error = null;
    _isLoading = true;
    notifyListeners();
    try {
      final result = await _authService.loginWithEmail(email, password);
      _token = result['token'];
      _user = result['user'];
      _isAuthenticated = true;
      _isLoading = false;
      notifyListeners();
      return true;
    } catch (e) {
      _error = e.toString().replaceAll('Exception: ', '');
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }

  Future<bool> register({
    required String name,
    required String email,
    required String password,
  }) async {
    _error = null;
    _isLoading = true;
    notifyListeners();
    try {
      final result = await _authService.register(
        name: name,
        email: email,
        password: password,
      );
      _token = result['token'];
      _user = result['user'];
      _isAuthenticated = true;
      _isLoading = false;
      notifyListeners();
      return true;
    } catch (e) {
      _error = e.toString().replaceAll('Exception: ', '');
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }

  // Social logins now pass user.id
  Future<bool> loginWithGoogle() async {
    _error = null;
    _isLoading = true;
    notifyListeners();
    try {
      final userId = _user?['id'] ?? 0;
      final result = await _authService.loginWithGoogle(userId: userId);
      _token = result['token'];
      _user = result['user'];
      _isAuthenticated = true;
      _isLoading = false;
      notifyListeners();
      return true;
    } catch (e) {
      _error = e.toString().replaceAll('Exception: ', '');
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }

  Future<bool> loginWithApple() async {
    _error = null;
    _isLoading = true;
    notifyListeners();
    try {
      final userId = _user?['id'] ?? 0;
      final result = await _authService.loginWithApple(userId: userId);
      _token = result['token'];
      _user = result['user'];
      _isAuthenticated = true;
      _isLoading = false;
      notifyListeners();
      return true;
    } catch (e) {
      _error = e.toString().replaceAll('Exception: ', '');
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }

  Future<bool> loginWithFacebook() async {
    _error = null;
    _isLoading = true;
    notifyListeners();
    try {
      final userId = _user?['id'] ?? 0;
      final result = await _authService.loginWithFacebook(userId: userId);
      _token = result['token'];
      _user = result['user'];
      _isAuthenticated = true;
      _isLoading = false;
      notifyListeners();
      return true;
    } catch (e) {
      _error = e.toString().replaceAll('Exception: ', '');
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }

  Future<void> logout() async {
    await _authService.logout();
    _token = null;
    _user = null;
    _isAuthenticated = false;
    _error = null;
    notifyListeners();
  }

  void clearError() {
    _error = null;
    notifyListeners();
  }
}