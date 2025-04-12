
        #include <windows.h>
        #include <FL/Fl.H>
        
        extern "C" void RunUI();
        
        int WINAPI WinMain(HINSTANCE hInstance, HINSTANCE hPrevInstance, LPSTR lpCmdLine, int nCmdShow) {
            RunUI();
            return 0;
        }
        
        int main(int argc, char** argv) {
            return WinMain(GetModuleHandle(NULL), NULL, GetCommandLine(), SW_SHOW);
        }
    