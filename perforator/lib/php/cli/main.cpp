#include <llvm/Support/TargetSelect.h>
#include <perforator/lib/php/php.h>
#include <perforator/lib/llvmex/llvm_exception.h>

#include <util/stream/format.h>

#include <llvm/Object/ObjectFile.h>

int main(int argc, const char* argv[]) {
    llvm::InitializeNativeTarget();
    llvm::InitializeNativeTargetDisassembler();

    Y_THROW_UNLESS(argc == 2);
    auto objectFile = Y_LLVM_RAISE(llvm::object::ObjectFile::createObjectFile(argv[1]));

    NPerforator::NLinguist::NPhp::TZendPhpAnalyzer analyzer{*objectFile.getBinary()};
    TMaybe<NPerforator::NLinguist::NPhp::TParsedPhpVersion> version = analyzer.ParseVersion();
    if (version) {
        Cout << "Parsed php binary version: "
             << version->ToString() << Endl;
    } else {
        Cout << "Could not parse php version" << Endl;
    }
    TMaybe<NPerforator::NLinguist::NPhp::EZendVmKind> vmKind = analyzer.ParseZendVmKind();
    if (vmKind) {
        Cout << "Parsed ZEND_VM_KIND: "
             << NPerforator::NLinguist::NPhp::ToString(*vmKind) << Endl;
    } else {
        Cout << "Could not parse ZEND_VM_KIND" << Endl;
    }

    TMaybe<bool> ztsEnabled = analyzer.ParseZts();
    if (ztsEnabled) {
        if (*ztsEnabled) {
            Cout << "ZTS enabled" << Endl;
        } else {
            Cout << "ZTS disabled" << Endl;
        }
    } else {
        Cout << "Could not parse ZTS" << Endl;
    }

    TMaybe<ui64> executorGlobalsAddress = analyzer.ParseExecutorGlobals();

    if (executorGlobalsAddress) {
        Cout << "Found executor_globals address: " << *executorGlobalsAddress << Endl;
    } else {
        Cout << "Could not find executor_globals" << Endl;
    }
}
