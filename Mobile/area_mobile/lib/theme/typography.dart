import 'package:flutter/material.dart';
import 'colors.dart';

/// Typography using Arial for mobile.
TextTheme buildTextTheme(AppColorPalette palette) {
  return TextTheme(
    displayLarge: TextStyle(
      fontSize: 32,
      fontWeight: FontWeight.w700,
      fontFamily: 'Arial',
      color: palette.almostBlack,
      letterSpacing: -0.5,
    ),
    titleLarge: TextStyle(
      fontSize: 22,
      fontWeight: FontWeight.w700,
      fontFamily: 'Arial',
      color: palette.deepBlue,
    ),
    titleMedium: TextStyle(
      fontSize: 18,
      fontWeight: FontWeight.w600,
      fontFamily: 'Arial',
      color: palette.deepBlue,
    ),
    bodyLarge: TextStyle(
      fontSize: 18,
      fontFamily: 'Arial',
      color: palette.almostBlack,
      height: 1.35,
    ),
    bodyMedium: TextStyle(
      fontSize: 16,
      fontFamily: 'Arial',
      color: palette.almostBlack,
      height: 1.4,
    ),
    bodySmall: TextStyle(
      fontSize: 14,
      fontFamily: 'Arial',
      color: palette.darkGrey,
    ),
    labelLarge: TextStyle(
      fontSize: 16,
      fontWeight: FontWeight.w600,
      fontFamily: 'Arial',
      color: palette.deepBlue,
      letterSpacing: 0.1,
    ),
  );
}
