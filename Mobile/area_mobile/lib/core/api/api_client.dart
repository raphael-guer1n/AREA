import 'package:dio/dio.dart';

class ApiClient {
  final Dio dio;

  ApiClient(String baseUrl)
      : dio = Dio(BaseOptions(baseUrl: baseUrl)) {
    dio.interceptors.add(LogInterceptor(responseBody: true));
  }

  Future<Response> get(String path) => dio.get(path);
}