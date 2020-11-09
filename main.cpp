#include <iostream>
#include <vector>
#include <openssl/bn.h>
#include <cstring>
#include <cmath>
#include <random>

using namespace std;

string PRIME = "115792089237316195423570985008687907853269984665640564039457584007913129640233"; //простое число для модульной арифметики
string ACCYRACY = "10000000000000000000000000000000000000000000000000000000000000000000000000000000000";
char SEP = ':';

struct Point {
    int64_t x;
    BIGNUM *y;
    Point (int64_t x, BIGNUM *y) {
        this->x=x;
        this->y=y;
    }
};


bool stringsToPoints(const vector<string>& str_points, vector <Point*> &points) { //перевод введённых строк в вектор структур Point

    for (auto str_point : str_points) {
        string str_tmp = str_point;

        int pos = str_tmp.find(SEP);
        if (pos<0 || pos >= str_point.length()-1) {
            printf("Ошибка: не удалось распознать введённые ключи\n");
            return false;
        }
        char *pEnd = nullptr;
        int64_t x = strtol(str_tmp.substr(pos+1).c_str(), &pEnd, 16);
        BIGNUM *y = nullptr;
        BN_hex2bn(&y, str_tmp.erase(pos).c_str());

        if (!x || y== nullptr) {
            printf("Ошибка: не удалось распознать введённые ключи\n%s\n", str_point.c_str());
            return false;
        }
        points.push_back(new Point(x, y));
    }
    return true;
}


BIGNUM *LagrangePolynomial(const vector<Point*>& points, double x=0) { //вычислить значение полинома по известным точкам для заданного х

    BN_CTX *ctx = BN_CTX_new();
    BIGNUM *res_coeff = nullptr;
    BN_dec2bn(&res_coeff, "0");
    BIGNUM *bn_coeff = nullptr; //множитель, необходимый для перехода к целочисленному делению
    BN_dec2bn(&bn_coeff, ACCYRACY.c_str());
    BIGNUM *bn_prime = nullptr;
    BN_dec2bn(&bn_prime, PRIME.c_str());

    BIGNUM *dividend = nullptr; //считаем числитель и значенатель отдельно
    BIGNUM *divisor = nullptr;

    for (int i=0; i<points.size(); i++) {
        BN_dec2bn(&dividend, "1");
        BN_dec2bn(&divisor, "1");

        BIGNUM *new_dividend = nullptr;
        BIGNUM *new_divisor = nullptr;

        for (int j=0; j<points.size(); j++)
            if (j!=i) {
                BN_dec2bn(&new_dividend, to_string(x-points[j]->x).c_str());
                BN_dec2bn(&new_divisor, to_string(points[i]->x-points[j]->x).c_str());

                BN_mul(dividend, dividend, new_dividend, ctx);
                BN_mul(divisor, divisor, new_divisor, ctx);
            }
        BN_mul(dividend, dividend, bn_coeff, ctx);

        BIGNUM *bn_mult_b = BN_new();
        BIGNUM *rem = BN_new();
        BIGNUM *round = BN_new();
        BN_div(bn_mult_b, rem, dividend, divisor, ctx); //целочисленное деление с округлением
        BN_add(rem, rem, rem);
        BN_div(round, rem, rem, divisor, ctx);
        BN_add(bn_mult_b, bn_mult_b, round);


        BIGNUM *new_add = BN_new();
        BN_mul(new_add, points[i]->y, bn_mult_b, ctx); //i-е слагаемое
        BN_add(res_coeff, res_coeff, new_add);

        BN_free(rem);
        BN_free(round);
        BN_free(bn_mult_b);
        BN_free(new_add);
        BN_free(new_dividend);
        BN_free(new_divisor);
    }
    BIGNUM *res = BN_new();
    BIGNUM *rem = BN_new();
    BIGNUM *round = BN_new();
    BN_div(res, rem, res_coeff, bn_coeff, ctx); //делим результат на домноженный коэффициент с учётом округления
    BN_add(rem, rem, rem);
    BN_div(round, rem, rem, bn_coeff, ctx);
    BN_add(res, res, round);


    BIGNUM *result_mod = BN_new();
    BN_mod(result_mod, res, bn_prime, ctx); //берём остаток от деления на PRIME

    BN_free(dividend);
    BN_free(divisor);
    BN_free(res_coeff);
    BN_free(bn_coeff);
    BN_free(bn_prime);
    BN_free(res);
    BN_free(rem);
    BN_free(round);

    return result_mod;
}

vector<BIGNUM*> generateCoefficients (int n, int64_t A=0, int64_t B = INT64_MAX) { //сгенерировать вектор случайных чисел
    random_device rd;
    mt19937 gen(rd());
    uniform_int_distribution<int64_t> dist (A, B);

    BN_CTX *ctx = BN_CTX_new();
    BIGNUM *pow = nullptr;
    BN_dec2bn(&pow, "4");

    vector<BIGNUM*> coeffs;
    for (int i=0; i<n; i++) {
        BIGNUM *new_coeff = nullptr;
        BN_dec2bn(&new_coeff, to_string(dist(gen)*dist(gen)*dist(gen)).c_str());
        BN_exp(new_coeff, new_coeff, pow, ctx);
        coeffs.push_back(new_coeff);
    }
    BN_free(pow);
    return coeffs;
}


BIGNUM* calculatePolynomial (vector<BIGNUM*> coeffs, int64_t x=0) { //вычислисть значение полинома по коэффициентам для заданного х

    BN_CTX *ctx = BN_CTX_new();
    BIGNUM *bn_x = nullptr;
    BN_dec2bn(&bn_x, to_string(x).c_str());

    BIGNUM *exp = BN_new();
    BIGNUM *mult = BN_new();
    BIGNUM *res = BN_new();
    BIGNUM *mod_res = BN_new();

    int n = coeffs.size();
    for (int i=0; i<n; i++) {

        BIGNUM *pow = nullptr; //res += x^(n-i-1) * coeffs[i]
        BN_dec2bn(&pow, to_string(n-i-1).c_str());

        BN_exp(exp, bn_x, pow, ctx);
        BN_mul(mult, exp, coeffs[i], ctx);
        BN_add(res, res, mult);

        BN_free(pow);
    }
    BIGNUM *bn_prime = nullptr;
    BN_dec2bn(&bn_prime, PRIME.c_str());
    BN_mod(mod_res, res, bn_prime, ctx); //берём остаток от деления результата на PRIME

    BN_free(exp);
    BN_free(mult);
    BN_free(res);
    BN_free(bn_x);
    BN_free(bn_prime);

    return mod_res;
}


void Recover() {

    vector <string> str_shares;
    string str;
    while (getline(cin, str) && (str.length()>0)) {
        str_shares.push_back(str);
    }
    vector <Point*> shares;
    if (!stringsToPoints(str_shares, shares)) return;

    BIGNUM *key = LagrangePolynomial(shares);
    char *str_key = BN_bn2hex(key); //Печать результата
    printf("%s\n", str_key);
    OPENSSL_free(str_key);
}

void Split() {

    int T, N;
    char KEY[256];
    cin >> N >> T;
    if ((N<T || N>99 || T<3)) {
        printf("Ошибка параметров! Необходимое условие: 2 < T <= N < 100");
        return;
    }
    cin >> KEY;
    BIGNUM *bn_KEY = nullptr;
    BN_hex2bn(&bn_KEY, KEY);

    vector<BIGNUM*> coeffs = generateCoefficients(T-1);
    coeffs.push_back(bn_KEY);

    for (int x=2; x<N+2; x++) {
        BIGNUM *new_y = calculatePolynomial(coeffs, x); //считаем значение полинома У для каждого х
        char *char_y = BN_bn2hex(new_y); //вывод результата
        printf("%s%c%d\n", char_y, SEP, x);
        OPENSSL_free(char_y);
    }
}


int main(int argc, char *argv[])
{
    if (argc!=2 || ((strcmp(argv[1], "split")!=0) && strcmp(argv[1], "recover")!=0)) {
        printf("Ошибка в указанных параметрах\n");
        return -1;
    }
    if (strcmp(argv[1], "split")==0)
        Split();
    else
        Recover();

    return 0;
}
