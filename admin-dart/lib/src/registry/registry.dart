class Registry {
  final String proto = "http";
  final String host = "localhost";
  final String port = "3000";

  String getFullUrl() {
    return proto + "://" + host + ":" + port;
  }
}