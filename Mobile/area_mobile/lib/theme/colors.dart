import 'package:flutter/material.dart';

enum VisionMode { normal, tritanopia }

/// Palette aligned with the Web globals.css tokens (light/dark + tritanopia)
class AppColorPalette extends ThemeExtension<AppColorPalette> {
  final Color white; // background
  final Color lightGrey; // surface / cards
  final Color grey; // borders
  final Color darkGrey; // muted text
  final Color almostBlack; // main text
  final Color navySoft; // accent shade 1
  final Color deepBlue; // primary
  final Color midBlue; // primary strong
  final Color softSlate; // accent / muted primary
  final Color placeholder;

  const AppColorPalette({
    required this.white,
    required this.lightGrey,
    required this.grey,
    required this.darkGrey,
    required this.almostBlack,
    required this.navySoft,
    required this.deepBlue,
    required this.midBlue,
    required this.softSlate,
    required this.placeholder,
  });

  factory AppColorPalette.forTheme({
    required Brightness brightness,
    required VisionMode vision,
  }) {
    final isDark = brightness == Brightness.dark;
    if (!isDark && vision == VisionMode.normal) {
      return const AppColorPalette(
        white: Color(0xFFFFFFFF),
        lightGrey: Color(0xFFF5F7FA),
        grey: Color(0xFFE2E8F0),
        darkGrey: Color(0xFF4A4A4A),
        almostBlack: Color(0xFF0D0D0D),
        navySoft: Color(0xFF0A1A2F),
        deepBlue: Color(0xFF112A46),
        midBlue: Color(0xFF1C3D63),
        softSlate: Color(0xFF41556B),
        placeholder: Color(0xFF9AA6B3),
      );
    }

    if (!isDark && vision == VisionMode.tritanopia) {
      return const AppColorPalette(
        white: Color(0xFFFFFFFF),
        lightGrey: Color(0xFFF5F5F9),
        grey: Color(0xFFE4E3ED),
        darkGrey: Color(0xFF4A4A4A),
        almostBlack: Color(0xFF0D0D0D),
        navySoft: Color(0xFF0A2525),
        deepBlue: Color(0xFF123938),
        midBlue: Color(0xFF1D5250),
        softSlate: Color(0xFF426160),
        placeholder: Color(0xFF9AA6B3),
      );
    }

    if (isDark && vision == VisionMode.normal) {
      return const AppColorPalette(
        white: Color(0xFF0A0F14),
        lightGrey: Color(0xFF0F1A24),
        grey: Color(0xFF1C2A36),
        darkGrey: Color(0xFFA8B4C0),
        almostBlack: Color(0xFFE8EDF2),
        navySoft: Color(0xFF0D253A),
        deepBlue: Color(0xFF163A59),
        midBlue: Color(0xFF1E4C73),
        softSlate: Color(0xFF7B8FA3),
        placeholder: Color(0xFF6F7F92),
      );
    }

    // Dark + tritanopia
    return const AppColorPalette(
      white: Color(0xFF0A1111),
      lightGrey: Color(0xFF0F1F1F),
      grey: Color(0xFF1C3030),
      darkGrey: Color(0xFFA8BABA),
      almostBlack: Color(0xFFE8EFEF),
      navySoft: Color(0xFF206260),
      deepBlue: Color(0xFF174B4A),
      midBlue: Color(0xFF0E3030),
      softSlate: Color(0xFF7C9A99),
      placeholder: Color(0xFF6F7F92),
    );
  }

  @override
  AppColorPalette copyWith({
    Color? white,
    Color? lightGrey,
    Color? grey,
    Color? darkGrey,
    Color? almostBlack,
    Color? navySoft,
    Color? deepBlue,
    Color? midBlue,
    Color? softSlate,
    Color? placeholder,
  }) {
    return AppColorPalette(
      white: white ?? this.white,
      lightGrey: lightGrey ?? this.lightGrey,
      grey: grey ?? this.grey,
      darkGrey: darkGrey ?? this.darkGrey,
      almostBlack: almostBlack ?? this.almostBlack,
      navySoft: navySoft ?? this.navySoft,
      deepBlue: deepBlue ?? this.deepBlue,
      midBlue: midBlue ?? this.midBlue,
      softSlate: softSlate ?? this.softSlate,
      placeholder: placeholder ?? this.placeholder,
    );
  }

  @override
  ThemeExtension<AppColorPalette> lerp(
    ThemeExtension<AppColorPalette>? other,
    double t,
  ) {
    if (other is! AppColorPalette) return this;
    return AppColorPalette(
      white: Color.lerp(white, other.white, t)!,
      lightGrey: Color.lerp(lightGrey, other.lightGrey, t)!,
      grey: Color.lerp(grey, other.grey, t)!,
      darkGrey: Color.lerp(darkGrey, other.darkGrey, t)!,
      almostBlack: Color.lerp(almostBlack, other.almostBlack, t)!,
      navySoft: Color.lerp(navySoft, other.navySoft, t)!,
      deepBlue: Color.lerp(deepBlue, other.deepBlue, t)!,
      midBlue: Color.lerp(midBlue, other.midBlue, t)!,
      softSlate: Color.lerp(softSlate, other.softSlate, t)!,
      placeholder: Color.lerp(placeholder, other.placeholder, t)!,
    );
  }
}

extension AppColorsX on BuildContext {
  AppColorPalette get appColors =>
      Theme.of(this).extension<AppColorPalette>() ??
      AppColorPalette.forTheme(
        brightness: Theme.of(this).brightness,
        vision: VisionMode.normal,
      );
}
