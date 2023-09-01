class ExecutionResult():

    def __init__(self, agency: int) -> None:
        self.agency = agency

    def ok(self) -> bool:
        return False
    
class OkResult(ExecutionResult):
    
    def __init__(self, agency: int=0) -> None:
        super().__init__(agency)

    def ok(self) -> bool:
        return True
    
class ErrResult(ExecutionResult):

    def __init__(self, agency: int=0) -> None:
        super().__init__(agency)
