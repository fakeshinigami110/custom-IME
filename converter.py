import re
from config_loader import load_mapping_config, load_engine_config, load_keywords_config


class BinaryConverter:
    def __init__(self):
        self.mapping = load_mapping_config()
        self.engine_config = load_engine_config()
        self.keywords, self.operators = load_keywords_config()

    def convert_text(self, text):
        """تبدیل متن به باینری"""
        if not text:
            return ""

        # ابتدا عملگرها را جایگزین کن اما با یک marker
        operator_replacements = {}
        for op, binary in sorted(self.operators.items(), key=lambda x: len(x[0]), reverse=True):
            if op in text:
                marker = f"__OP_{len(operator_replacements)}__"
                operator_replacements[marker] = binary
                text = text.replace(op, f" {marker} ")

        # سپس کلمات کلیدی را جایگزین کن
        words = text.split()
        result = []

        for word in words:
            # اگر marker عملگر است
            if word.startswith("__OP_") and word in operator_replacements:
                result.append(operator_replacements[word])
            # اعداد (شامل 0 و 1)
            elif self._is_number(word):
                result.append('$' + self._number_to_binary(word))
            # کلمات کلیدی
            elif word in self.keywords:
                result.append(self.keywords[word])
            # ترکیب حروف و اعداد
            elif self._has_letters_and_numbers(word):
                result.append(self._convert_mixed_simple(word))
            # حروف و کلمات
            else:
                result.append(self._convert_letters(word))

        return ' '.join(result)

    def _is_number(self, word):
        """بررسی آیا کلمه عدد است"""
        if word.startswith("__OP_"):
            return False

        if word.startswith('-'):
            rest = word[1:]
        else:
            rest = word

        if rest.replace('.', '').isdigit() and rest.count('.') <= 1:
            return True
        return False

    def _has_letters_and_numbers(self, word):
        """بررسی آیا کلمه ترکیب حروف و اعداد است"""
        if word.startswith("__OP_"):
            return False

        has_letters = any(c.isalpha() for c in word)
        has_numbers = any(c.isdigit() for c in word)
        return has_letters and has_numbers

    def _convert_mixed_simple(self, text):
        """تبدیل ترکیب حروف و اعداد - روش ساده"""
        result = ""
        i = 0
        n = len(text)

        while i < n:
            # اگر کاراکتر عددی است
            if text[i].isdigit():
                # پیدا کردن کل عدد
                j = i
                while j < n and (text[j].isdigit() or text[j] == '.'):
                    j += 1
                number_part = text[i:j]
                result += "$" + self._number_to_binary(number_part)
                i = j
            else:
                # پیدا کردن کل حروف
                j = i
                while j < n and text[j].isalpha():
                    j += 1
                letter_part = text[i:j]
                result += self._convert_letters(letter_part)
                i = j

        return result

    def _convert_letters(self, text):
        """تبدیل حروف به باینری"""
        result = ""
        for char in text:
            if char.lower() in self.mapping:
                result += self.mapping[char.lower()]
            else:
                result += char
        return result

    def _number_to_binary(self, number_str):
        """تبدیل عدد به باینری"""
        try:
            is_negative = number_str.startswith('-')
            clean_number = number_str[1:] if is_negative else number_str

            if '.' in clean_number:
                integer_part, fractional_part = clean_number.split('.')
                integer_bin = bin(int(integer_part))[2:] if integer_part else "0"
                fractional_bin = self._fraction_to_binary(fractional_part)
                result = f"{integer_bin}.{fractional_bin}"
            else:
                result = bin(int(clean_number))[2:]

            return f"0 {result}" if is_negative else result

        except (ValueError, Exception):
            return number_str

    def _fraction_to_binary(self, fractional_str, precision=8):
        """تبدیل بخش اعشاری به باینری"""
        try:
            if not fractional_str or fractional_str == "0":
                return "0"

            fractional = int(fractional_str) / (10 ** len(fractional_str))
            result = ""

            for _ in range(precision):
                fractional *= 2
                if fractional >= 1:
                    result += "1"
                    fractional -= 1
                else:
                    result += "0"
                if fractional == 0:
                    break

            return result if result else "0"
        except (ValueError, Exception):
            return fractional_str


# تست
if __name__ == "__main__":
    converter = BinaryConverter()

    print("=== Final Fixed Tests ===")
    tests = [
        "45hello",  # عدد سپس حروف
        "hello42world",  # حروف سپس عدد سپس حروف
        "test123",  # حروف سپس عدد
        "123test",  # عدد سپس حروف
        "a20b30",  # حروف-عدد-حروف-عدد
        "0",  # عدد صفر
        "1",  # عدد یک
        "x <= 12 + 1",  # اعداد با عملگرها
    ]

    for test in tests:
        result = converter.convert_text(test)
        print(f"'{test}': '{result}'")