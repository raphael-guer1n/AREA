import 'package:flutter/material.dart';
import '../theme/colors.dart';
import 'home/home_screen.dart';
import 'area/area_screen.dart';
import 'services/services_screen.dart';
import 'profile/profile_screen.dart';

class MainShell extends StatefulWidget {
  const MainShell({super.key});

  @override
  State<MainShell> createState() => _MainShellState();
}

class _MainShellState extends State<MainShell> {
  int _selectedIndex = 0;

  // List of pages in the order they appear in the nav bar
  final List<Widget> _pages = const [
    HomeScreen(),
    AreaScreen(),
    ServicesScreen(),
    ProfileScreen(),
  ];

  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Scaffold(
      backgroundColor: colorScheme.background,
      body: SafeArea(child: _pages[_selectedIndex]),

      // AREA Design System inspired navigation bar
      bottomNavigationBar: Container(
        decoration: BoxDecoration(
          color: colorScheme.surface,
          border: const Border(
            top: BorderSide(color: AppColors.grey, width: 0.6),
          ),
          boxShadow: [
            BoxShadow(
              color: AppColors.grey.withOpacity(0.15),
              blurRadius: 6,
              offset: const Offset(0, -2),
            ),
          ],
        ),
        child: BottomNavigationBar(
          type: BottomNavigationBarType.fixed,
          currentIndex: _selectedIndex,
          onTap: _onItemTapped,
          backgroundColor: colorScheme.surface,
          selectedItemColor: colorScheme.primary,
          unselectedItemColor: AppColors.darkGrey.withOpacity(0.7),
          selectedLabelStyle: theme.textTheme.labelLarge,
          unselectedLabelStyle: theme.textTheme.bodySmall,
          elevation: 0,
          items: const [
            BottomNavigationBarItem(
              icon: Icon(Icons.dashboard_rounded),
              label: 'Dashboard',
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.sync_alt_rounded),
              label: 'AREA',
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.widgets_outlined),
              label: 'Services',
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.person_outline_rounded),
              label: 'Profile',
            ),
          ],
        ),
      ),
    );
  }
}