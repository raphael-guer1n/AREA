import "package:flutter/foundation.dart";
import "../services/auth_service.dart";

class AuthProvider extends ChangeNotifier {
  final AuthService _authService = AuthService();

  bool _isAuthenticated = false;
  bool _isLoading = false;
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
      if (savedToken != null && savedToken.isNotEmpty) {
        _token = savedToken;

        try {
          final fetchedUser = await _authService.fetchCurrentUser();
          _user = fetchedUser;
          _isAuthenticated = true;
        } catch (e) {
          debugPrint("Failed to fetch current user: $e");
        }
      }
    } catch (e) {
      debugPrint("Error checking auth status: $e");
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
      _token = result["token"] as String?;
      _user = result["user"] as Map<String, dynamic>?;
      _isAuthenticated = true;

      _isLoading = false;
      notifyListeners();
      return true;
    } catch (e) {
      _error = e.toString().replaceAll("Exception: ", "");
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
      _token = result["token"] as String?;
      _user = result["user"] as Map<String, dynamic>?;
      _isAuthenticated = true;

      _isLoading = false;
      notifyListeners();
      return true;
    } catch (e) {
      _error = e.toString().replaceAll("Exception: ", "");
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }

  Future<bool> loginWithGoogleForLogin() async {
    _error = null;
    _isLoading = true;
    notifyListeners();

    try {
      final result = await _authService.loginWithGoogleWithoutUser();
      _token = result['token'];

      final userData = await _authService.fetchCurrentUser();
      _user = userData;
      _isAuthenticated = true;

      _isLoading = false;
      notifyListeners();
      return true;
    } catch (e) {
      _error = e.toString().replaceAll("Exception: ", "");
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
}