import 'package:flutter/material.dart';
import '../services/config_service.dart';
import '../theme/colors.dart';
import '../theme/theme.dart';

class ThemeProvider extends ChangeNotifier {
  VisionMode _visionMode = VisionMode.normal;
  bool _loaded = false;

  ThemeProvider() {
    _loadVisionMode();
  }

  VisionMode get visionMode => _visionMode;
  bool get isLoaded => _loaded;

  ThemeData get lightTheme =>
      buildAreaTheme(brightness: Brightness.light, vision: _visionMode);
  ThemeData get darkTheme =>
      buildAreaTheme(brightness: Brightness.dark, vision: _visionMode);

  Future<void> setVisionMode(VisionMode mode) async {
    if (_visionMode == mode) return;
    _visionMode = mode;
    notifyListeners();
    await ConfigService.setVisionMode(mode.name);
  }

  Future<void> _loadVisionMode() async {
    final stored = await ConfigService.getVisionMode();
    if (stored == VisionMode.tritanopia.name) {
      _visionMode = VisionMode.tritanopia;
    }
    _loaded = true;
    notifyListeners();
  }
}
