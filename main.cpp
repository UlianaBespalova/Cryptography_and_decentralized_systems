#include <iostream>
#include <fstream>
#include <vector>
#include <cmath>
#include <cstring>

using namespace std;


vector<string> getStringsFromFile (const string& fileName) {
    vector<string> lines;
    string newLine;
    ifstream fin (fileName);

    if (!fin.is_open()) {
        cout << "Не удалось открыть файл" << endl;
        return lines;
    }
    while (getline(fin, newLine))
        lines.push_back(newLine);
    return lines;
}

void divideNumber(int N, int &left, int &right) {
    int num = 1;
    int tmpN = N;
    while ((tmpN/=10)>0) num++;

    int tmp = pow(10, (num/2));
    left = N/tmp;
    right = N%tmp;
}

int convertToNumber(const string& line, int N, int sid=2) {
    int sum = 1;
    for (auto sym: line) {
        sum = (sum * sid + sym);
    }
    int numLeft, numRight;
    divideNumber(abs(sum), numLeft, numRight);
    int result = ((unsigned)numLeft^(unsigned)numRight) % N + 1;

    return result;
}



int main(int argc, char *argv[]) {

    int N, param;
    string fileName;
    for (int i=1; i<argc; i++) {
        char *pEnd;
        if (i+1!=argc) {
            if (strcmp(argv[i], "-f")==0) {
                fileName = argv[i+1];
                i++;
            }
            else if (strcmp(argv[i], "-n")==0) {
                N = (int)strtol(argv[i+1], &pEnd, 10);
                i++;
            }
            else if (strcmp(argv[i], "-p")==0) {
                param = (int)strtol(argv[i+1], &pEnd, 10);
                i++;
            }
            else cout << "Неверный ключ - " << i << endl;
        }
    }
    if (N<=0) {
        cout << "Ошибка: N должно быть положительным" << endl;
        return 0;
    }



    vector<string> lines = getStringsFromFile(fileName);

    for(auto & line : lines) {
        int number = convertToNumber(line, N, param);
        cout << line << " : " << number << endl;
    }
    return 0;
}