import os
import tempfile
import unittest

from config_parser.implementation.env_parser import EnvParser
from config_parser.implementation.ini_parser import IniParser


class TestEnvParser(unittest.TestCase):
    def test_basic_key_value(self) -> None:
        p = EnvParser("KEY=value")
        self.assertEqual(p["KEY"], "value")

    def test_multiple_lines(self) -> None:
        source = "A=1\nB=2\nC=3"
        p = EnvParser(source)
        self.assertEqual(p["A"], "1")
        self.assertEqual(p["B"], "2")
        self.assertEqual(p["C"], "3")

    def test_comments_and_empty_lines(self) -> None:
        source = "# comment\n\nKEY=val\n\n# another"
        p = EnvParser(source)
        self.assertEqual(p["KEY"], "val")

    def test_quoted_values(self) -> None:
        source = 'KEY="value with spaces"'
        p = EnvParser(source)
        self.assertEqual(p["KEY"], "value with spaces")

    def test_single_quoted_values(self) -> None:
        source = "KEY='value with spaces'"
        p = EnvParser(source)
        self.assertEqual(p["KEY"], "value with spaces")

    def test_escaped_chars(self) -> None:
        source = 'KEY="line1\\nline2"'
        p = EnvParser(source)
        self.assertIn("\n", p["KEY"])

    def test_variable_substitution(self) -> None:
        source = "HOME=/home/user\nPATH=$HOME/bin"
        p = EnvParser(source)
        self.assertEqual(p["PATH"], "/home/user/bin")

    def test_variable_substitution_braces(self) -> None:
        source = "HOME=/home/user\nPATH=${HOME}/bin"
        p = EnvParser(source)
        self.assertEqual(p["PATH"], "/home/user/bin")

    def test_env_var_substitution(self) -> None:
        with unittest.mock.patch.dict(os.environ, {"TEST_VAR": "env_value"}):
            source = "KEY=$TEST_VAR"
            p = EnvParser(source)
            self.assertEqual(p["KEY"], "env_value")

    def test_parse_file(self) -> None:
        with tempfile.NamedTemporaryFile(mode="w", suffix=".env", delete=False) as f:
            f.write("DB_HOST=localhost\nDB_PORT=5432\n")
            f.flush()
            p = EnvParser()
            p.parse_file(f.name)
            self.assertEqual(p["DB_HOST"], "localhost")
            self.assertEqual(p["DB_PORT"], "5432")
            os.unlink(f.name)

    def test_missing_key_returns_default(self) -> None:
        p = EnvParser("A=1")
        self.assertIsNone(p.get("B"))
        self.assertEqual(p.get("B", "default"), "default")

    def test_contains(self) -> None:
        p = EnvParser("KEY=val")
        self.assertIn("KEY", p)
        self.assertNotIn("NOPE", p)


class TestIniParser(unittest.TestCase):
    def test_basic_section(self) -> None:
        source = "[section]\nkey=value"
        p = IniParser(source)
        self.assertEqual(p.get("section", "key"), "value")

    def test_multiple_sections(self) -> None:
        source = "[db]\nhost=localhost\n[app]\nport=8080"
        p = IniParser(source)
        self.assertEqual(p.get("db", "host"), "localhost")
        self.assertEqual(p.get("app", "port"), "8080")

    def test_comments(self) -> None:
        source = "; comment\n# also comment\n[key]\nval=test"
        p = IniParser(source)
        self.assertEqual(p.get("key", "val"), "test")

    def test_default_section(self) -> None:
        source = "default_key=default_val\n[section]\nkey=value"
        p = IniParser(source)
        self.assertEqual(p.get("section", "default_key"), "default_val")

    def test_quoted_values(self) -> None:
        source = '[sec]\nval="quoted"'
        p = IniParser(source)
        self.assertEqual(p.get("sec", "val"), "quoted")

    def test_type_coercion_int(self) -> None:
        source = "[sec]\ncount=42"
        p = IniParser(source)
        self.assertEqual(p.getint("sec", "count"), 42)

    def test_type_coercion_float(self) -> None:
        source = "[sec]\npi=3.14"
        p = IniParser(source)
        self.assertEqual(p.getfloat("sec", "pi"), 3.14)

    def test_type_coercion_bool(self) -> None:
        source = "[sec]\na=true\nb=false\nc=yes\nd=no"
        p = IniParser(source)
        self.assertTrue(p.getbool("sec", "a"))
        self.assertFalse(p.getbool("sec", "b"))
        self.assertTrue(p.getbool("sec", "c"))
        self.assertFalse(p.getbool("sec", "d"))

    def test_type_coercion_defaults(self) -> None:
        source = "[sec]\na=1"
        p = IniParser(source)
        self.assertIsNone(p.getint("sec", "nonexistent"))
        self.assertEqual(p.getint("sec", "nonexistent", 99), 99)

    def test_sections_list(self) -> None:
        source = "[a]\nk=v\n[b]\nk=v\n"
        p = IniParser(source)
        self.assertCountEqual(p.sections(), ["a", "b"])

    def test_missing_section(self) -> None:
        source = "[a]\nk=v"
        p = IniParser(source)
        with self.assertRaises(KeyError):
            _ = p["nonexistent"]

    def test_getitem(self) -> None:
        source = "[sec]\nk1=v1\nk2=v2"
        p = IniParser(source)
        self.assertEqual(p["sec"], {"k1": "v1", "k2": "v2"})

    def test_parse_file(self) -> None:
        with tempfile.NamedTemporaryFile(mode="w", suffix=".ini", delete=False) as f:
            f.write("[section]\nkey=value\n")
            f.flush()
            p = IniParser()
            p.parse_file(f.name)
            self.assertEqual(p.get("section", "key"), "value")
            os.unlink(f.name)

    def test_colon_separator(self) -> None:
        source = "[sec]\nkey: value"
        p = IniParser(source)
        self.assertEqual(p.get("sec", "key"), "value")

    def test_type_coercion_error(self) -> None:
        source = "[sec]\nval=not_a_number"
        p = IniParser(source)
        with self.assertRaises(ValueError):
            p.getint("sec", "val")
        with self.assertRaises(ValueError):
            p.getbool("sec", "val")


if __name__ == "__main__":
    unittest.main()
