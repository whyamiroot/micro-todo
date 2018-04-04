import 'package:angular/angular.dart';
import 'registry.dart';
import 'package:admin_dart/src/service/service.dart';

@Injectable()
class RegistryService {
  final Registry registry = new Registry();

  List<String> getListOfServiceTypes() {
    return new List<String>();
  }

  List<Service> getListOfServicesByType(String type) {
    return new List<Service>();
  }
}