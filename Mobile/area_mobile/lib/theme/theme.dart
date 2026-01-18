import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'colors.dart';
import 'typography.dart';

/// Builds ThemeData aligned with Web globals.css (light/dark + tritanopia).
ThemeData buildAreaTheme({
  required Brightness brightness,
  required VisionMode vision,
}) {
  final palette =
      AppColorPalette.forTheme(brightness: brightness, vision: vision);

  final colorScheme = ColorScheme(
    brightness: brightness,
    primary: palette.deepBlue,
    onPrimary: Colors.white,
    secondary: palette.softSlate,
    onSecondary: Colors.white,
    error: Colors.red.shade600,
    onError: Colors.white,
    background: palette.white,
    onBackground: palette.almostBlack,
    surface: palette.lightGrey,
    onSurface: palette.almostBlack,
  );

  return ThemeData(
    useMaterial3: true,
    brightness: brightness,
    scaffoldBackgroundColor: palette.white,
    primaryColor: palette.deepBlue,
    colorScheme: colorScheme,
    textTheme: buildTextTheme(palette),
    dividerTheme: DividerThemeData(
      color: palette.grey.withOpacity(0.7),
      thickness: 1,
    ),
    cardTheme: CardThemeData(
      color: palette.white,
      surfaceTintColor: palette.white,
      elevation: 2,
      margin: const EdgeInsets.symmetric(vertical: 8, horizontal: 12),
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.all(Radius.circular(12)),
      ),
      shadowColor: palette.grey.withOpacity(0.3),
    ),
    inputDecorationTheme: InputDecorationTheme(
      fillColor: palette.lightGrey,
      filled: true,
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(10),
        borderSide: BorderSide(color: palette.grey),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(10),
        borderSide: BorderSide(color: palette.deepBlue, width: 1.5),
      ),
      labelStyle: TextStyle(color: palette.darkGrey),
      hintStyle: TextStyle(color: palette.placeholder),
      prefixIconColor: palette.darkGrey,
      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
    ),
    elevatedButtonTheme: ElevatedButtonThemeData(
      style: ElevatedButton.styleFrom(
        backgroundColor: palette.midBlue,
        foregroundColor: Colors.white,
        textStyle: GoogleFonts.josefinSans(
          fontWeight: FontWeight.w700,
          fontSize: 16,
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(10),
        ),
        padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
      ),
    ),
    outlinedButtonTheme: OutlinedButtonThemeData(
      style: OutlinedButton.styleFrom(
        foregroundColor: palette.deepBlue,
        side: BorderSide(color: palette.grey),
        textStyle: GoogleFonts.josefinSans(
          fontWeight: FontWeight.w600,
          fontSize: 15,
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(10),
        ),
      ),
    ),
    bottomNavigationBarTheme: BottomNavigationBarThemeData(
      backgroundColor: palette.white,
      selectedItemColor: palette.deepBlue,
      unselectedItemColor: palette.darkGrey.withOpacity(0.7),
      selectedIconTheme: IconThemeData(color: palette.deepBlue),
      unselectedIconTheme:
          IconThemeData(color: palette.darkGrey.withOpacity(0.7)),
      elevation: 8,
    ),
    appBarTheme: AppBarTheme(
      backgroundColor: Colors.transparent,
      elevation: 0,
      foregroundColor: palette.almostBlack,
      centerTitle: false,
      titleTextStyle: GoogleFonts.josefinSans(
        color: palette.almostBlack,
        fontWeight: FontWeight.w700,
        fontSize: 20,
      ),
      iconTheme: IconThemeData(color: palette.deepBlue),
    ),
    snackBarTheme: SnackBarThemeData(
      behavior: SnackBarBehavior.floating,
      backgroundColor: palette.deepBlue,
      contentTextStyle: GoogleFonts.josefinSans(color: Colors.white),
    ),
    extensions: <ThemeExtension<dynamic>>[
      palette,
    ],
  );
}
