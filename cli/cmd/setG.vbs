Set wshShell = CreateObject("WScript.Shell")
Set env = wshShell.Environment("User")
old = env("Path")
Dim currentDir
currentDir = CreateObject("Scripting.FileSystemObject").GetAbsolutePathName(".")
newVar = currentDir & "\run-scripts"
if InStr(1, old, newVar, vbTextCompare) = 0 Then
	old = newVar & ";" & old
End If
env("Path") = old
