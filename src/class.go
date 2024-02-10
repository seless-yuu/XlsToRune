package main

import (
	"os"
	"path"
	"strconv"
)

func OutputClassString(book RuneTypeBook, json_path string, out_dir string) error {
	for _, sheet := range book.Sheets {
		err := outputSheet(sheet, json_path, book.Name, out_dir)
		if err != nil {
			return err
		}
	}

	return nil
}

func outputSheet(sheet RuneTypeSheet, file_name string, json_path string, out_dir string) error {
	for _, table := range sheet.Tables {
		err := outputTable(table, json_path, file_name, out_dir)
		if err != nil {
			return err
		}
	}

	return nil
}

func outputTable(table RuneTypeTable, file_name string, json_path string, out_dir string) error {
	class_name := file_name + "_" + table.Name
	json_dir := path.Dir(json_path)

	class_str := "using System;\n"
	class_str += "using UnityEngine;\n"
	class_str += "using UnityEngine.AddressableAssets;\n"
	class_str += "using UnityEngine.ResourceManagement.AsyncOperations;\n"
	class_str += "using RuneImporter;\n"
	class_str += "\n"

	class_str += "namespace RuneImporter\n"
	class_str += "{\n"
	class_str += "    public static partial class RuneLoader\n"
	class_str += "    {\n"
	class_str += "        public static AsyncOperationHandle " + class_name + "_LoadInstanceAsync()\n"
	class_str += "        {\n"
	class_str += "            return " + "Rune." + class_name + ".LoadInstanceAsync();\n"
	class_str += "        }\n"
	class_str += "    }\n"
	class_str += "}\n"
	class_str += "\n"

	class_str += "namespace Rune\n"
	class_str += "{\n"
	class_str += "\n"

	class_str += addRuneClassName(class_name, len(table.Values))

	for _, t := range table.Types {
		switch t.TypeName.Kind {
		case SString:
			class_str += addRuneString(t.TypeName)
		case SInt:
			class_str += addRuneInt(t.TypeName)
		case SInt2:
			class_str += addRuneInt2(t.TypeName)
		case SInt3:
			class_str += addRuneInt3(t.TypeName)
		case SFloat:
			class_str += addRuneFloat(t.TypeName)
		case SFloat2:
			class_str += addRuneFloat2(t.TypeName)
		case SFloat3:
			class_str += addRuneFloat3(t.TypeName)
		case SFloat4:
			class_str += addRuneFloat4(t.TypeName)
		}
	}

	class_str += "    }\n"

	class_str += "\n"
	class_str += "    public static AsyncOperationHandle LoadInstanceAsync() {\n"
	class_str += "        var src_dir = \"" + json_dir + "/\";\n"
	class_str += "        var out_dir = string.IsNullOrEmpty(Config.ScriptableObjectDirectory) ? src_dir : Config.ScriptableObjectDirectory;\n"
	class_str += "        var asset_name = \"" + table.Name + ".asset\";\n"
	class_str += "        var path = out_dir + asset_name;\n"
	class_str += "        var handle = Config.OnLoad(path);\n"
	class_str += "        handle.Completed += (handle) => { instance = handle.Result as " + class_name + "; };\n"
	class_str += "\n"
	class_str += "        return handle;\n"
	class_str += "    }\n"

	class_str += "}\n"
	class_str += "\n"
	class_str += "}\n"

	out_file_name := "Rune_" + file_name + "_" + table.Name
	return write(out_file_name, class_str, out_dir)
}

func write(class_name string, class_str string, out_dir string) error {
	path := out_dir + "/" + class_name + ".cs"

	err := os.MkdirAll(out_dir, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte(class_str))
	if err != nil {
		return err
	}

	return nil
}

func addRuneClassName(type_name string, value_length int) string {
	str := "public class " + type_name + " : RuneScriptableObject\n"
	str += "{\n"
	str += "    public static " + type_name + " instance { get; private set; }\n"
	str += "\n"
	str += "    [SerializeField]\n"
	str += "    public Value[] ValueList = new Value[" + strconv.Itoa(value_length) + "];\n"
	str += "\n"
	str += "    [Serializable]\n"
	str += "    public struct Value\n"
	str += "    {\n"

	return str
}

func addRuneString(type_name RuneTypeName) string {
	return "        public string " + type_name.Value + ";\n"
}

func addRuneInt(type_name RuneTypeName) string {
	return "        public int " + type_name.Value + ";\n"
}

func addRuneInt2(type_name RuneTypeName) string {
	return "        public Int2 " + type_name.Value + ";\n"
}

func addRuneInt3(type_name RuneTypeName) string {
	return "        public Int3 " + type_name.Value + ";\n"
}

func addRuneFloat(type_name RuneTypeName) string {
	return "        public float " + type_name.Value + ";\n"
}

func addRuneFloat2(type_name RuneTypeName) string {
	return "        public Vector2 " + type_name.Value + ";\n"
}

func addRuneFloat3(type_name RuneTypeName) string {
	return "        public Vector3 " + type_name.Value + ";\n"
}

func addRuneFloat4(type_name RuneTypeName) string {
	return "        public Vector4 " + type_name.Value + ";\n"
}
